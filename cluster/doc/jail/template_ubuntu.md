# NAME

**template-ubuntu** - create a Ubuntu jail template


# HOST

Enable
[Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/#linuxemu)
to run Linux user land on FreeBSD:

```console
# sysrc -f /etc/rc.conf.d/linux linux_enable=YES
# service linux start
```

The service loads kernel modules and mounts file systems for Linux applications
under `/compat/linux`:

```console
% ls /compat/linux
dev proc  sys
```

Linux application can run on the host like native FreeBSD binaries now. They
are present in the process tree, can be traced, etc.

FreeBSD Ports tree includes a small number of Linux applications with `linux-`
prefix:

```console
% pkg search '^linux-*' | head -n 3
linux-ai-ml-env-1.0.0          Linux Python environment for running Stable Diffusion models and PyTorch CUDA examples
linux-bcompare-4.4.7           Compare, sync, and merge files and folders (X11)
linux-bitwarden-cli-1.22.1     Bitwarden CLI
```

Beware that when run, these applications expect to find configurations,
libraries, and other files under `/` root folder. Use `chroot(8)` to guarantee
its working:

```console
# chroot /compat/linux /bin/ls
```

# JAIL

## Debian bootstrap

Create a FreeBSD jail template `Ubuntu-22.04`:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p2 zroot/jail/template/Ubuntu-22.04
```

Manually start a jail at the template to install bootstrap Linux shared
libraries inside the template:

```console
# cat /tmp/jammy.conf
jammy {
  host.hostname = "${name}.lab.net";
  path = '/jail/template/Ubuntu-24.04';
  exec.consolelog = "/var/log/jail_${name}.log";

  interface = "em0";
  ip4.addr = "192.168.1.10";

  allow.raw_sockets;
  exec.clean;

  mount.devfs;
  devfs_ruleset = 4;

  allow.mount;
  allow.mount.devfs;
  allow.mount.fdescfs;
  allow.mount.procfs;
  allow.mount.linprocfs;
  allow.mount.linsysfs;
  allow.mount.tmpfs;

  enforce_statfs = 1;

  exec.start = "/bin/sh /etc/rc";
  exec.stop = "/bin/sh /etc/rc.shutdown";
}
# jail -cm -f /tmp/jammy.conf
```

Verify network setup:

```console
# jexec jammy ifconfig em0 | grep inet
  inet 192.168.1.10 netmask 0xffffffff broadcast 192.168.1.10
# jexec jammy ping -c 1 192.168.1.1
PING 192.168.1.1 (192.168.1.1): 56 data bytes
64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=1.446 ms

--- 192.168.1.1 ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 1.446/1.446/1.446/0.000 ms
# jexec jammy nc -uv 1.1.1.1 53
Connection to 1.1.1.1 53 port [udp/domain] succeeded!
^C
```

Use
[debootstrap(8)](https://manpages.debian.org/stretch/debootstrap/debootstrap.8.en.html)
to provide Linux shared libraries.

```console
# jexec jammy pkg install debootstrap
The package management tool is not yet installed on your system.
Do you want to fetch and install it now? [y/N]: y
...
```

<details>
<summary>Package messages</summary>

```
=====
Message from debootstrap-1.0.128n2_4:

--
To successfully create an installation of Debian or Ubuntu
debootstrap requires the following kernel modules to be loaded:

linux64 fdescfs linprocfs linsysfs tmpfs

To install Ubuntu 18.04 LTS (Bionic Beaver) into /compat/ubuntu, run as root:

debootstrap bionic /compat/ubuntu
```

</details>

`debootstrap`:

  * It is a shell script from Debian - https://wiki.debian.org/Debootstrap.
  * FreeBSD 14.3 uses
    [1.0.128 version](https://cgit.freebsd.org/ports/commit/sysutils/debootstrap/Makefile?id=2bf4a73e61cd322efa426f55101afa25bd2481d3)
    as of Oct '25.
  * `debootstrap` uses scripts from `/usr/local/share/debootstrap/scripts/`
    named after Linux distribution.
  * Version 1.0.128 does not include Ubuntu 24.04 Noble but has Ubuntu 22.04
    Jammy.
  * Install Linux libraries into `/compat/<distribution>` to emphasize the
    distribution in use.

```console
# jexec jammy debootstrap jammy /compat/jammy
I: Retrieving InRelease
I: Checking Release signature
...
I: Base system installed successfully.
```

Stop the jail:

```console
# jail -r -f /tmp/jammy.conf
```

Snapshot the template:

```console
# zfs snapshot zroot/jail/template/Ubuntu-22.04@p1
```

Clean up temporary jail configuration:

```console
# rm /etc/jail.conf.d/jammy.conf
```
