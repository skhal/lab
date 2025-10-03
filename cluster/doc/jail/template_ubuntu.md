# NAME

**template-ubuntu** - create a Ubuntu jail template


# DESCRIPTION

## Host

Enable
[Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/#linuxemu)
to run Linux user land on FreeBSD:

```console
# sysrc -f /etc/rc.conf.d/linux linux_enable=YES
# service linux start
```

The service
[loads kernel modules](https://github.com/freebsd/freebsd-src/blob/ae5914c0e4478fd35ef9db3f32665b60e04d5a6f/libexec/rc/rc.d/linux#L32-L63)
and
[mounts file systems](https://github.com/freebsd/freebsd-src/blob/ae5914c0e4478fd35ef9db3f32665b60e04d5a6f/libexec/rc/rc.d/linux#L74-L80)
for Linux applications under `/compat/linux` prefix:

```console
# sysctl -n compat.linux.emul_path
/compat/linux
```

Now Linux applications can run on the host along native FreeBSD binaries now.
They show up in the process tree, can be traced, etc. There is a small number
of such applications available in the FreeBSD Ports tree with `linux-` prefix:

```console
% pkg search '^linux-*' | head -n 3
linux-ai-ml-env-1.0.0          Linux Python environment for running Stable Diffusion models and PyTorch CUDA examples
linux-bcompare-4.4.7           Compare, sync, and merge files and folders (X11)
linux-bitwarden-cli-1.22.1     Bitwarden CLI
```

## Jail

debootstrap(8) installs Linux shared libraries. We'll need to install these
inside the template. The template also needs to access default FreeBSD setup.

### ZFS Bootstrap

Start from FreeBSD template:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p2 zroot/jail/template/Ubuntu-22.04
```

### Start template

Manually start a jail with path set to the template:

```console
# cat /tmp/jammy.conf
jammy {
  host.hostname = "${name}.lab.net";
  path = "/jail/template/Ubuntu-22.04";

  interface = "em0";
  ip4.addr = "192.168.1.10";

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
# jail -cm -f /tmp/jammy.conf
jammy: created
```

Verify work:

```console
# jexec jammy ifconfig em0
em0: flags=1008943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	options=a520b9<RXCSUM,VLAN_MTU,VLAN_HWTAGGING,JUMBO_MTU,VLAN_HWCSUM,WOL_MAGIC,VLAN_HWFILTER,VLAN_HWTSO,RXCSUM_IPV6,HWSTATS>
	ether 00:1f:c6:9c:54:f1
	inet 192.168.1.10 netmask 0xffffffff broadcast 192.168.1.10
	media: Ethernet autoselect (1000baseT <full-duplex>)
	status: active
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

### Debian bootstrap

Install debootstrap(8):

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

> [!IMPORTANT]
> The message from debootstrap(8) includes installed version `1.0.128`.
>
> It is a Debian shell script, https://wiki.debian.org/Debootstrap. The script
> pulls a number of libraries, customized for different Debian distributions by
> `/usr/local/share/debootstrap/scripts/`.
>
> For example, `jammy` script is for Ubuntu 22.04.
>
> debootstrap(8) v1.0.128 does not include Ubuntu 24.04 Noble.

Install Linux libraries into `/compat/<distribution>`, i.e., `/compat/jammy`:

```console
# jexec jammy debootstrap jammy /compat/jammy
I: Retrieving InRelease
...
I: Base system installed successfully.
```

Verify work:

```console
# jexec jammy ls /compat/jammy
bin dev home  lib32 libx32  mnt proc  run srv tmp var
boot  etc lib lib64 media opt root  sbin  sys usr
```

### Stop template

Stop the jail:

```console
# jail -r -f /tmp/jammy.conf jammy
jammy: removed
```

### Snapshot

Create a ZFS Snapshot the template. Name it after FreeBSD patch. Include the
change version:

```console
# zfs snapshot zroot/jail/template/Ubuntu-22.04@p2.0
```
