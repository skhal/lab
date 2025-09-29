# NAME

**jail** - setup jails on `nuc.lab.net`


# DESCRIPTION

Ref: https://docs.freebsd.org/en/books/handbook/jails/

  * Jail configurations are under `/etc/jail.conf.d/`
  * Jail filesystem is under ZFS dataset `zroot/jail`

# HOST SETUP

## Jail service

Start jails in parallel, reverse the shutdown order and enable `jail` service.

```console
# sysrc -f /etc/rc.conf.d/jail jail_parallel_start="YES"
jail_parallel_start: NO -> YES
# sysrc -f /etc/rc.conf.d/jail jail_reverse_stop="YES"
jail_reverse_stop: NO -> YES
# sysrc -f /etc/rc.conf.d/jail jail_enable="YES"
jail_enable: NO -> YES
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

https://github.com/skhal/lab/blob/23207daeb5f8ed9708d40eb8141ef3c1157a744c/cluster/net.lab.nuc/etc/jail.conf#L16


# CREATE TEMPLATE

Download userland (`-m` forces to download the file only if the remote version
is newer than the local file):

```console
# fetch https://download.freebsd.org/ftp/releases/amd64/amd64/14.3-RELEASE/base.txz -m -o /jail/image/14.3-RELEASE-base.txz
```

Create a template with updated software:

> [!TIP]
> Set `PAGER` to `cat(1)` to suppress interactive mode for `freebsd-update(8)`.

```console
# tar -xf /jail/image/14.3-RELEASE-base.txz -C /jail/template/14.3-RELEASE --unlink
# env PAGER=cat freebsd-update -b /jail/template/14.3-RELEASE/ fetch install
# env ROOT=/jail/template/14.3-RELEASE /jail/template/14.3-RELEASE/bin/freebsd-version -u
14.3-RELEASE-p2
```

Fix timezone from UTC to local (US Central Time):

```console
# chroot /jail/template/14.3-RELEASE date
Tue Sep 16 15:45:41 UTC 2025
# chroot /jail/template/14.3-RELEASE tzsetup
# chroot /jail/template/14.3-RELEASE date
Tue Sep 16 10:47:15 CDT 2025
```

Change root shell to `tcsh(1)`:

```console
# chroot /jail/template/14.3-RELEASE chsh -s /bin/tcsh root
```

Copy root tcsh(1) resource configuration:

```console
# cp /root/.cshrc /jail/template/14.3-RELEASE/root
```

Operate syslogd(8) without sockets, i.e., local mode:

```console
# chroot /jail/template/14.3-RELEASE sysrc -f /etc/rc.conf.d/syslogd syslogd_flags="-ss"
```

Set default gateway to avoid `routing` service restarts:

```console
# chroot /jail/template/14.3-RELEASE sysrc -f /etc/rc.conf.d/routing defaultrouter="192.168.1.1"
```

Snapshot the version with installed patch number.

```console
# zfs snapshot zroot/jail/template/14.3-RELEASE@p1
```

> [!NOTE]
> Check latest snapshot before using it with
> `zfs list -t snapshot zroot/jail/template/14.3-RELEASE`

Include [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail) for jail VNET management:

```console
# fetch -o /usr/local/etc/rc.jail https://raw.githubusercontent.com/skhal/lab/refs/heads/main/freebsd/rc/rc.jail
```


# NO NET JAIL

The new jail `dev` is isolated and does not have a network access:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p1 zroot/jail/container/dev
# cat /etc/jail.conf.d/dev.conf
dev {
  # empty
}
```

# VNET JAIL

Ref: https://freebsdfoundation.org/wp-content/uploads/2020/03/Jail-vnet-by-Examples.pdf

Virtual Network (VNET) adds a networking stack to the jail, isolated from the
host system. It includes interfaces, addresses, routing tables and firewall
rules.

The setup uses two network interfaces to connect jails using bridge(4):

* `bridge0` connects jails to the Internet. Subnet: `192.168.1.0/24`
* `bridge1` isolates intra-jail traffic. Subnet: `10.0.1.0/24`

Create bridges and give `bridge0` Internet access:

```console
# sysrc -f /etc/rc.conf.d/network cloned_interfaces+="bridge0 bridge1"
cloned_interfaces:  -> bridge0 bridge1
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge0="addm em0 up descr jail:em"
ifconfig_bridge0:  -> addm em0 up descr jail:em
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge1="up descr jail:lo"
ifconfig_bridge0:  -> addm em0 up descr jail:lo
```

> [!Note]
> The bridge and both ends of epair must be in the UP state for the packets to
> flow.

Let [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail)
script managet jail's epair(4), connect it to the bridge, and assign the IP
address.

```console
% cat /etc/jail.conf.d/dev.conf
dev {
  $id = "110";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  vnet;
  devfs_ruleset = 5;

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```

# LINUX JAIL

Ref: https://docs.freebsd.org/en/books/handbook/jails/#creating-linux-jail

## Step 1: Linux Binary Compatibility

Enable [Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/#linuxemu)
on the host node `nuc.lab.net`:

```console
# sysrc -f /etc/rc.conf.d/linux linux_enable=YES
# service linux start
```

The service loads kernel modules and mounts filesystems for Linux applications
under /compat/linux:

```console
% ls /compat/linux
dev proc  sys
```

> [!NOTE]
> Linux applications can be started like native FreeBSD binaries, behave like a
> native process, can be traced and debugged in the same way. There is a number
> of Linux applications available through Ports tree with `linux-` prefix:
>
> ```console
> % pkg search '^linux-*' | head -n 3
> linux-ai-ml-env-1.0.0          Linux Python environment for running Stable Diffusion models and PyTorch CUDA examples
> linux-bcompare-4.4.7           Compare, sync, and merge files and folders (X11)
> linux-bitwarden-cli-1.22.1     Bitwarden CLI
> ```

## Step 2: Linux Userland in VNET Jail

Create a `ubuntu` VNET jail with `$id=111`.

Provide Linux shared libraries using
[debootstrap(8)](https://manpages.debian.org/stretch/debootstrap/debootstrap.8.en.html).

```console
# pkg install debootstrap
# debootstrap jammy /compat/ubuntu
```

where `jammy` is the [Ubuntu release](https://www.releases.ubuntu.com) name,
i.e. LTS Ubuntu 22.04 Jammy Jellyfish. Install Linux libraries into
`/compat/<distribution>` instead of `/compat/linux` to make it clear what
distribution is in use.

## Step 3: Linux Mount Points

At this point, the host environemnt has Linux ABI enabled by `linux` service,
which [mounts](https://github.com/freebsd/freebsd-src/blob/1c3ca0c733a4e4ba550cedfa8019260fb0cf5707/libexec/rc/rc.d/linux#L75-L79)
a number of filesystems for Linux in the host environment: `devfs`, `procfs`,
etc.

We'll need to share these filesystems with the jail for Ubuntu applications to
operate correctly.

> [!WARNING]
> Mount filesystems under `$path/compat/ubuntu` since we'll run Linux
> applications with `chroot /compat/ubuntu ...`.

Stop the jail and add the following instructions to the jail configuration:

```
# service jail stop ubuntu
# cat /etc/jail.conf.d/ubuntu
ubuntu {
  $id = "111";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  allow.mount;
  allow.mount.devfs;
  allow.mount.fdescfs;
  allow.mount.procfs;
  allow.mount.linprocfs;
  allow.mount.linsysfs;
  allow.mount.tmpfs;

  enforce_statfs = 1; # only mount points below jail's chroot

  vnet;
  devfs_ruleset = 5;

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";

  mount += "devfs     $path/compat/ubuntu/dev     devfs     rw  0 0";
  mount += "tmpfs     $path/compat/ubuntu/dev/shm tmpfs     rw,size=1g,mode=1777  0 0";
  mount += "fdescfs   $path/compat/ubuntu/dev/fd  fdescfs   rw,linrdlnk 0 0";
  mount += "linprocfs $path/compat/ubuntu/proc    linprocfs rw  0 0";
  mount += "linsysfs  $path/compat/ubuntu/sys     linsysfs  rw  0 0";
}
```

> [!NOTE]
> FreeBSD instructions include `/tmp` and `/home` mounts. Mount these only if
> using X11 applications.

Start the jail:

```console
# service jail start ubuntu
```

Verify Linux userland:

```console
# jexec ubuntu
root@ubuntu:/ # chroot /compat/ubuntu uname -s -r -m
Linux 5.15.0 x86_64
```

## Finish: Summary

The new jail `ubuntu` is a standard FreeBSD userland with Linux ABI enabled and
Linux libraries included under `/compat/ubuntu`. The jail does not have
`/compat/linux` because Linux ABI runs in the host environment.

If we run any of the Linux software inside the jail from under
`/compat/ubuntu/bin`, it will user Linux ABI and Linux libraries from
`/compat/ubuntu`.

Linux applications run in the same space with FreeBSD applications inside the
jail. One can check running applications with standard FreeBSD tools:

```console
root@ubuntu:/ # ps -Adf
  PID TT  STAT    TIME COMMAND
 9526  -  IsJ  0:00.00 /usr/sbin/sshd
13221  -  IsJ  0:00.01 /usr/sbin/cron -s
97816  -  SsJ  0:00.01 /usr/sbin/syslogd -s
51490  1  SJ   0:00.04 /bin/tcsh -i
96496  1  R+J  0:00.00 - ps -Adf
```

The jail runs services from different environments: FreeBSD `systlogd` and
Linux `sshd`:

```console
root@ubuntu:/ # service sshd status
sshd is not running.
root@ubuntu:/ # service syslogd status
syslogd is running as pid 97816.
root@ubuntu:/ # chroot /compat/ubuntu /usr/sbin/service --status-all
 [ - ]  console-setup.sh
 [ - ]  cron
 [ - ]  dbus
 [ ? ]  hwclock.sh
 [ - ]  keyboard-setup.sh
 [ ? ]  kmod
 [ - ]  procps
 [ + ]  ssh
 [ - ]  udev
 [ - ]  unattended-upgrades
```

Keep in mind that Linux applications expect Linux filesystem hierarchy. Use
`chroot(8)`:

```console
root@ubuntu:/ # /compat/ubuntu/usr/bin/uname -a
ELF interpreter /lib64/ld-linux-x86-64.so.2 not found, error 2
Abort
root@ubuntu:/ # chroot /compat/ubuntu /usr/bin/uname -a
Linux ubuntu.nuc.lab.net 5.15.0 FreeBSD 14.3-RELEASE-p2 GENERIC x86_64 x86_64 x86_64 GNU/Linux
```

Therefore it is helpful to enter the jail with `chroot(8)` unless one needs to
debug FreeBSD jail:

```console
# jexec ubuntu chroot /compat/ubuntu su -l
```

or get into the jail the usual way:

```console
# jexec ubuntu su -l
```
