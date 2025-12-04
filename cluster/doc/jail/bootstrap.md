NAME
====

**host** - host setup for jails

DESCRIPTION
===========

ZFS
---

Jails reside under `/jail`. There are three folders:

-	`/jail/image` stores downloaded user lands in compressed format.
-	`/jail/template` holds base templates to create jails.
-	`/jail/container` keeps running jail.

Create datasets:

```console
# zfs create -o mountpoint=/jail zroot/jail
# zfs create zroot/jail/image
# zfs create zroot/jail/template
# zfs create zroot/jail/container
```

Permissions
-----------

Create `jail` user group:

```console
# pw groupadd -g 1001 -n jail
```

`jail` members can manage and enter jails:

```console
# cat <<eof >> /usr/local/etc/doas.conf
permit nopass :jail cmd jail
permit nopass :jail cmd jexec
eof
```

Members of `jail` can manage jail datasets:

```console
# zfs allow -s @mount mount,canmount,mountpoint zroot/jail
# zfs allow -s @create create,destroy,@mount zroot/jail
# zfs allow -g jail @mount,@create,readonly,snapshot zroot/jail
```

Add system operator to the `jail` group:

```console
# pw groupmod -m op -n jail
```

Jail service
------------

`jail(8)` reads configuration from `/etc/jail.conf`. Set default configuration parameters and pick up individual jail configurations from `/etc/jail.conf.d/`

```console
# mkdir /etc/jail.conf.d
# cat <<eof >/etc/jail.conf
.include "/etc/jail.conf.d/*.conf";
eof
```

Enable `jail` service and stop jails in the reverse order to ensure dependencies are satisfied:

```console
# sysrc -f /etc/rc.conf.d/jail jail_reverse_stop=yes
jail_reverse_stop: NO -> yes
# sysrc -f /etc/rc.conf.d/jail jail_enable=yes
jail_enable: NO -> yes
```

Network
-------

Jails use vnet(9), virtual network, to isolate networking setup from the host environment. Use if_bridge(4) to connect networks - network interfaces that are UP and connected to the bridge in the UP state can pass packets.

Setup two bridges to separate local and Internet traffic.

Create a bridge for Internet traffic:

```console
# sysrc -f /etc/rc.conf.d/network cloned_interfaces+=bridge0
cloned_interfaces: '' -> bridge0
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge0='addm em0 up descr jail:em'
ifconfig_bridge0: '' -> addm em0 up descr jail:em
```

There are two Internet access markets in the bridge:

1.	it has em(4) Ethernet adapter connected.
2.	it has `jail:em` description to indicate that the bridge is for jails with em(4) adapter.

Create a second bridge for local, intra-jail traffic:

```console
# sysrc -f /etc/rc.conf.d/network cloned_interfaces+=bridge1
cloned_interfaces: '' -> bridge0
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge1='up descr jail:lo'
ifconfig_bridge0: '' -> up descr jail:lo
```

Again, it has `jail:lo` marker in the description to indicate that the bridge is for jails local traffic.

For each jail, we'll create an epair(4), a virtual back-to-back connected Ethernet interface, with a-side connected to the bridge in the host environment, and b-side moved to the jail environment, where it is assigned an IP address and brought UP.

We'll use [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail) script to manage epair(4) setup and tear down on per-bridge, per-jail basis.

```console
# fetch \
    -o /usr/local/etc/rc.jail \
    https://raw.githubusercontent.com/skhal/lab/refs/heads/main/freebsd/rc/rc.jail
```

SEE ALSO
========

-	https://docs.freebsd.org/en/books/handbook/jails/
