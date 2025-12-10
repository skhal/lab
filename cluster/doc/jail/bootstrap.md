Name
====

**bootstrap** - bootstrap FreeBSD to host jails

ZFS
===

Jails file hierarchy under `/jail`:

-	`/jail/image` stores downloaded userlands in compressed format.
-	`/jail/template` keeps base templates to create jails.
-	`/jail/container` holds running jails.

Create datasets:

```console
# zfs create -o mountpoint=/jail zroot/jail
# zfs create zroot/jail/image
# zfs create zroot/jail/template
# zfs create zroot/jail/container
```

Permissions
===========

Create `jail` user group to let members of the group manage jails:

```console
# pw groupadd -g 1001 -n jail
# pw groupmod -n jail -m op
```

`jail` members can manage and enter jails:

```console
# cat <<EOF >> /usr/local/etc/doas.conf
permit nopass :jail cmd jail
permit nopass :jail cmd jexec
EOF
```

Members of `jail` can snapshot jails:

```console
# zfs allow -g jail snapshot zroot/jail
```

Jail service
------------

`jail(8)` reads configuration from `/etc/jail.conf`. Set default configuration parameters and pick up individual jail configurations from `/etc/jail.conf.d/`

```console
# mkdir /etc/jail.conf.d
# cat <<EOF >/etc/jail.conf
.include "/etc/jail.conf.d/*.conf";
EOF
```

Enable `jail` service and stop jails in the reverse order to ensure dependencies are satisfied:

```console
% doas sysrc -f /etc/rc.conf.d/jail jail_reverse_stop=yes
jail_reverse_stop: NO -> yes
% doas sysrc -f /etc/rc.conf.d/jail jail_enable=yes
jail_enable: NO -> yes
```

Network
-------

Jails use vnet(9), virtual network, to isolate networking setup from the host environment. Use if_bridge(4) to connect networks - network interfaces that are UP and connected to the bridge in the UP state can pass packets.

Setup two bridges to separate local and Internet traffic.

Create a bridge for Internet traffic:

```console
% doas sysrc -f /etc/rc.conf.d/network cloned_interfaces+=bridge0
cloned_interfaces: '' -> bridge0
% doas sysrc -f /etc/rc.conf.d/network ifconfig_bridge0='addm igc0 up descr jail:ext'
ifconfig_bridge0: '' -> addm igc0 up descr jail:ext
```

The bridge has following properties:

1.	it has igc(4) Ethernet adapter connected for external connections.
2.	it has `jail:ext` description to indicate that the bridge is with external connection.

Create a second bridge for local, intra-jail traffic:

```console
% doas sysrc -f /etc/rc.conf.d/network cloned_interfaces+=bridge1
cloned_interfaces: 'bridge0' -> bridge0 bridge1
% doas sysrc -f /etc/rc.conf.d/network ifconfig_bridge1='inet 10.0.1.101/24 up descr jail:int'
ifconfig_bridge1: '' -> up descr jail:int
```

The new bridge has the following properties:

1.	it has an IP assigned.
2.	there are no connected Ethernet adapters.
3.	it has `jail:int` marker in the description to indicate that the bridge is for intra-jail traffic.

Create an epair(4) for each jail, a virtual back-to-back connected Ethernet interface, with two sides:

-	the a-side connects to the bridge in the host environment and remain in the hosting environment.
-	the b-side moves to the jail environment. The jail manages this side of the interface. It assigns an IP address, brings it UP, etc.

Use [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail) script to manage epair(4) setup and tear down on per-bridge, per-jail basis.

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/freebsd/rc/rc.jail
# mv -nv /tmp/rc.jail /usr/local/etc/rc.jail
```

References
==========

-	https://docs.freebsd.org/en/books/handbook/jails/
