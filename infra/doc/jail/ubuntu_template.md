<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**ubuntu-template** - create a Ubuntu jail template

Bootstrap
=========

Host
----

Enable [Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/#linuxemu) to run Linux software in FreeBSD environment, controlled by `linux` service.

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
# zfs clone zroot/jail/template/15.0-RELEASE@p0.0 zroot/jail/template/Ubuntu-22.04
```

We'll use a temporary jail to configure the template, with configuration in `/tmp/ubuntu.conf`. It is important to start and stop the jail using jail(8) command with `-f /tmp/ubuntu.conf` flag.

Start a temporary jail with path set to the template location:

```console
% fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/jail.conf.d/ubuntu_bootstrap.conf
% : update the path in /tmp/ubuntu_bootstrap.conf
% doas jail -cm -f /tmp/ubuntu_bootstrap.conf
ubuntu: created
```

Verify work:

```console
% doas jexec ubuntu ping -c 1 192.168.1.1
PING 192.168.1.1 (192.168.1.1): 56 data bytes
64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=0.585 ms

--- 192.168.1.1 ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 0.585/0.585/0.585/0.000 ms
% doas jexec ubuntu nc -uv 1.1.1.1 53
Connection to 1.1.1.1 53 port [udp/domain] succeeded!
^C
```

Use latest packages:

```console
% cat /jail/template/Ubuntu-22.04/usr/local/etc/pkg/repos/FreeBSD.conf
FreeBSD-ports: {
      url: "pkg+https://pkg.FreeBSD.org/${ABI}/latest",
}
% doas jexec ubuntu pkg update
```

Install debootstrap(8) to create Linux environment at a given prefix with minimal set of Linux shared libraries and binaries for Linux user land, including apt(1).

```console
% doas jexec ubuntu pkg install debootstrap
The package management tool is not yet installed on your system.
Do you want to fetch and install it now? [y/N]: y
...
```

<details>
<summary>Package messages</summary>

```
% doas jexec ubuntu pkg info -D debootstrap
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
% doas jexec ubuntu debootstrap jammy /compat/jammy
I: Retrieving InRelease
...
I: Base system installed successfully.
```

Verify work:

```console
% doas jexec ubuntu ls /compat/jammy
bin dev home  lib32 libx32  mnt proc  run srv tmp var
boot  etc lib lib64 media opt root  sbin  sys usr
```

Stop the jail:

```console
% doas jail -r -f /tmp/ubuntu_bootstrap.conf ubuntu
ubuntu: removed
```

Configure
=========

Some of the configuration steps such as apt(8) need access to mountpoints like `/dev`. Modify the jail's temporary configuration to include mounts:

```console
% fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/jail.conf.d/ubuntu_config.conf
% : update the path in /tmp/ubuntu_config.conf
% doas jail -cm -f /tmp/ubuntu_config.conf
ubuntu: created
```

Apt
---

Use Universe and Multiverse sources to apt(8)

```console
% doas jexec ubuntu fetch -o /compat/jammy/etc/apt/ https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/sources.list
```

Upgrade user land:

```console
% doas jexec ubuntu chroot /compat/jammy apt update
% doas jexec ubuntu chroot /compat/jammy apt upgrade
```

Install following packages for basic FreeBSD-like tools:

```console
% doas jexec ubuntu chroot /compat/jammy apt install apt-file ldnsutils man net-tools tcsh vim
```

-	`apt-file` to search for packages that install a file, i.e., `apt-file find /usr/bin/shasum`.
-	`net-tools` for ifconfig(1).
-	`ldnsutils` for drill(1).

Locale
------

```console
% doas jexec ubuntu chroot /compat/jammy locale-gen C.UTF-8
Generating locales (this might take a while)...
  C.UTF-8... done
Generation complete.
% doas jexec ubuntu chroot /compat/jammy dpkg-reconfigure locales
```

Root
----

Change shell to tcsh(1):

```console
% doas jexec ubuntu chroot /compat/jammy chsh -s /usr/bin/tcsh root
```

Copy `.cshrc`:

```console
# cp /root/.cshrc /jail/template/Ubuntu-22.04/root/
# cp /root/.cshrc /jail/template/Ubuntu-22.04/compat/jammy/root/
```

Set password:

```console
% doas jexec ubuntu chroot /compat/jammy passwd
```

Timezone
--------

```console
% doas jexec ubuntu chroot /compat/jammy dpkg-reconfigure tzdata
% doas jexec ubuntu chroot /compat/jammy date
Tue Nov 18 08:49:04 CST 2025
```

Snapshot
========

Stop the jail:

```console
% doas jail -r -f /tmp/ubuntu_config.conf ubuntu
ubuntu: removed
```

Snapshot the template:

```console
% zfs snapshot zroot/jail/template/Ubuntu-22.04@p0.0
```

SEE ALSO
========

-	https://docs.freebsd.org/en/books/handbook/jails/#creating-linux-jail
