# NAME

**jail** - setup jails on `nuc.lab.net`


# DESCRIPTION

Ref: https://docs.freebsd.org/en/books/handbook/jails/

  * Jail configurations are under `/etc/jail.conf.d/`
  * Jail filesystem is under ZFS dataset `zroot/jail` 

# HOST SETUP

## Jail service

Stop jails in the reverse order to resolve jail dependencies:

```console
# sysrc -f /etc/rc.conf.d/jail jail_reverse_stop="YES"
jail_reverse_stop: NO -> YES
```

## User groups

Create a jail group to manipulate jails:

```console
# pw groupadd -g 1001 -n jail
# pw groupmod -m op -n jail
```

## Doas

Let members of `jail` group control jails:

```console
# cat <<eof > /usr/local/etc/doas.conf
permit nopass :jail cmd jail
permit nopass :jail cmd jexec
eof
```

## ZFS datasets

Create ZFS root dataset for jails:

```console
# zfs create -o mountpoint=/jail zroot/jail
```

Grant group `jail` permissions to control datasets:

```control
# zfs allow -s @mount  mount,canmount,mountpoint zroot/jail
# zfs allow -s @create create,destroy,@mount zroot/jail
# zfs allow -g jail @mount,@create,readonly zroot/jail
```

Create child datasets:

  * `zroot/jail/image` stores loaded images
  * `zroot/jail/template` holds base templates for jails
  * `zroot/jail/container` running jails

```
# zfs create zroot/jail/image
# zfs create zroot/jail/template
# zfs create zroot/jail/container
```

## Jail config

Include jail configuraitons from `/etc/jail.conf.d/`:

```console
# cat <<eof >/etc/jail.conf
.include "/etc/jail.conf.d/*.conf";
eof
```


# CREATE TEMPLATE

Download userland (`-m` forces to download the file only if the remote version
is newer than the local file):

```console
# fetch https://download.freebsd.org/ftp/releases/amd64/amd64/14.3-RELEASE/base.txz -m -o /jail/image/14.3-RELEASE-base.txz
```

Create a template with updated software:

```console
# tar -xf /jail/image/14.3-RELEASE-base.txz -C /jail/template/14.3-RELEASE --unlink
# env PAGER=cat freebsd-update -b /jail/template/14.3-RELEASE/ fetch install
# env ROOT=/jail/template/14.3-RELEASE /jail/template/14.3-RELEASE/bin/freebsd-version -u
14.3-RELEASE-p2
# zfs snapshot zroot/jail/template/14.3-RELEASE@p2
```

> [!NOTE]
> Set `PAGER` to `cat(1)` to suppress interactive mode for `freebsd-update(8)`.


# NO NET JAIL

The new jail `dev` is isolated and does not have a network access:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p2 zroot/jail/container/dev
# cat <<eof >/etc/jail.conf.d/dev.conf
dev {
  exec.start = "/bin/sh /etc/rc";
  exec.stop = "/bin/sh /etc/rc.shutdown";
  exec.consolelog = "/var/log/jail_console_${name}.log";

  host.hostname = "${name}.nuc.lab.net";
  path = "/jail/container/${name}";

  exec.clean;
}
eof
```

# VNET JAIL

Virtual Network (VNET) adds a networking stack to the jail, isolated from the
host system. It includes interfaces, addresses, routing tables and firewall
rules.

Ref: https://freebsdfoundation.org/wp-content/uploads/2020/03/Jail-vnet-by-Examples.pdf

TL;DR: create an `epair(4)` on the host system, jail one end of it and connect
the other end to one of the network interfaces on the host system via
`bridge(4)`.

> [!Note]
> The bridge and both ends of epair must be in the UP state for the packets to 
> flow.

Create a bridge on the host system and add the network interface with Internet
access, `em0`:

```console
# sysrc -f /etc/rc.conf.d/network cloned_interfaces+="bridge0"
cloned_interfaces:  -> bridge0
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge0="addm em0 up"
ifconfig_bridge0:  -> addm em0 up
```

The jail configuration:

  * *Host system*: create an epair and add the a-side of it to the bridge in
    `exec.prestart`. Remove the epair from the bridge and destroy it in
    `exec.poststop`.

  * *Jail*: enable VNET and enajil the b-side of the epair into jail (it will
    be auto-released when the jail shuts down). Assign an IP address to the
    enjailed network interface.

```console
% cat /etc/jail.conf.d/dev.conf
dev {
  $id = "110";
  $ip = "192.168.1.${id}/24";
  $epair = "epair${id}";
  $bridge = "bridge0";

  allow.raw_sockets;
  exec.clean;
  mount.devfs;
  devfs_ruleset = 5;

  host.hostname = "${name}.nuc.lab.net";
  path = "/jail/container/${name}";

  # Enable VNET and enjail the b-side of the epair
  vnet;
  vnet.interface = "${epair}b";
  
  exec.consolelog = "/var/log/jail_console_${name}.log";
  
  exec.prestart  = "cp /etc/resolv.conf ${path}/etc/";
  exec.prestart += "/sbin/ifconfig ${epair} create";
  exec.prestart += "/sbin/ifconfig ${epair}a up";
  exec.prestart += "/sbin/ifconfig ${bridge} addm ${epair}a up";

  exec.start  = "/sbin/ifconfig ${epair}b ${ip} up";
  exec.start += "/bin/sh /etc/rc";

  exec.stop = "/bin/sh /etc/rc.shutdown";

  exec.poststop += "/sbin/ifconfig ${bridge} deletem ${epair}a";
  exec.poststop += "/sbin/ifconfig ${epair}a destroy";
}
```
