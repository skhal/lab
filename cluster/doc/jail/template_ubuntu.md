NAME
====

**template-ubuntu** - create a Ubuntu jail template

DESCRIPTION
===========

The instructions enable Linux Binary Compatibility on the hosting machine, create a ZFS dataset with bootstrapped Ubuntu user land.

Background
----------

FreeBSD provides [Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/#linuxemu) to run Linux software in FreeBSD environment, controlled by `linux` service.

debootstrap(8) creates Linux environment at a given prefix with minimal set of Linux shared libraries and binaries for Linux user land, including apt(1).

Bootstrap
---------

Enable Linux Binary Compatibility on the FreeBSD host:

```console
# sysrc -f /etc/rc.conf.d/linux linux_enable=YES
# service linux start
```

The service [loads kernel modules](https://github.com/freebsd/freebsd-src/blob/ae5914c0e4478fd35ef9db3f32665b60e04d5a6f/libexec/rc/rc.d/linux#L32-L63) and [mounts file systems](https://github.com/freebsd/freebsd-src/blob/ae5914c0e4478fd35ef9db3f32665b60e04d5a6f/libexec/rc/rc.d/linux#L74-L80) for Linux applications under `/compat/linux` prefix:

```console
% sysctl -n compat.linux.emul_path
/compat/linux
```

There is a limited number of Linux applications available in pkg(1) with `linux-` prefix:

```console
% pkg search '^linux-' | wc -l
     279
```

Keep in mind that Linux binaries run along FreeBSD binaries. They show up in the process tree, can be traced, etc.

Jail
----

Create a Ubuntu template jail from FreeBSD template:

```console
# zfs list -t snapshot zroot/jail/template/14.3-RELEASE | tail -1
zroot/jail/template/14.3-RELEASE@p5.6     0B      -   459M  -
# zfs clone zroot/jail/template/14.3-RELEASE@p5.6 zroot/jail/template/Ubuntu-22.04
```

Start a temporary jail with path set to the template location:

```console
% cat /tmp/jammy.conf
jammy {
  host.hostname = "${name}.lab.net";
  path = "/jail/template/Ubuntu-22.04";

  interface = "em0";
  ip4.addr = "192.168.1.12";

  allow.mount.devfs;
  allow.mount.fdescfs;
  allow.mount.linprocfs;
  allow.mount.linsysfs;
  allow.mount.procfs;
  allow.mount.tmpfs;
  allow.mount;
  allow.raw_sockets;

  mount.devfs;
  devfs_ruleset = 4;

  enforce_statfs = 1;

  exec.clean;
  exec.consolelog = "/var/log/jail_${name}.log";

  exec.start = "/bin/sh /etc/rc";
  exec.stop = "/bin/sh /etc/rc.shutdown";
}
% doas jail -cm -f /tmp/jammy.conf
jammy: created
```

Verify work:

```console
% doas jexec jammy ifconfig -ag epair
epair6b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:jammy:bridge0
	options=8<VLAN_MTU>
	ether 02:55:ef:72:1f:0b
	inet 192.168.1.12 netmask 0xffffff00 broadcast 192.168.1.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
epair7b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:jammy:bridge1
	options=8<VLAN_MTU>
	ether 02:54:f8:a8:4a:0b
	inet 10.0.1.12 netmask 0xffffff00 broadcast 10.0.1.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
% doas jexec jammy ping -c 1 192.168.1.1
PING 192.168.1.1 (192.168.1.1): 56 data bytes
64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=0.585 ms

--- 192.168.1.1 ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 0.585/0.585/0.585/0.000 ms
% doas jexec jammy nc -uv 1.1.1.1 53
Connection to 1.1.1.1 53 port [udp/domain] succeeded!
^C
```

Install debootstrap(8):

```console
% doas jexec jammy pkg install debootstrap
The package management tool is not yet installed on your system.
Do you want to fetch and install it now? [y/N]: y
...
```

<details>
<summary>Package messages</summary>

```
% doas jexec jammy pkg info -D debootstrap
debootstrap-1.0.128n2_4:
On install:
To successfully create an installation of Debian or Ubuntu
debootstrap requires the following kernel modules to be loaded:

linux64 fdescfs linprocfs linsysfs tmpfs

To install Ubuntu 18.04 LTS (Bionic Beaver) into /compat/ubuntu, run as root:

debootstrap bionic /compat/ubuntu
```

</details>

> [!IMPORTANT] debootstrap(8) is a Debian shell script, https://wiki.debian.org/Debootstrap. It pulls a bunch of libraries and apps, customized for different Debian distributions by `/usr/local/share/debootstrap/scripts/`. For example, `jammy` script aims to reproduce Ubuntu 22.04 distribution.
>
> debootstrap(8) v1.0.128 does not include a script for Ubuntu 24.04 Noble.

Bootstrap Linux land at `/compat/<distribution>`:

```console
% doas jexec jammy debootstrap jammy /compat/jammy
I: Retrieving InRelease
...
I: Base system installed successfully.
```

Verify work:

```console
% doas jexec jammy ls /compat/jammy
bin dev home  lib32 libx32  mnt proc  run srv tmp var
boot  etc lib lib64 media opt root  sbin  sys usr
```

Stop the jail:

```console
% doas jail -r -f /tmp/jammy.conf jammy
jammy: removed
```

Snapshot the template (include installed FreeBSD patch in the name):

```console
# zfs snapshot zroot/jail/template/Ubuntu-22.04@p5.0
```

SEE ALSO
========

-	https://docs.freebsd.org/en/books/handbook/jails/#creating-linux-jail
