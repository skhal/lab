# NAME

**jail-ubuntu** - create a Ubuntu jail


# DESCRIPTION

Clone Ubuntu jail template:

```console
# zfs clone zroot/jail/template/Ubuntu-22.04@p1 zroot/jail/container/jammy
```

## Mount points

The host environment runs `linux` service. It
[mounts](https://github.com/freebsd/freebsd-src/blob/ad38f6a0b466bf05a0d40ce1daa8c7bce0936271/libexec/rc/rc.d/linux#L75-L79)
different file systems for Linux under:

```console
# sysctl -n compat.linux.emul_path
/compat/linux
```

Propagate these mount points into Linux jail under `/compat/<distribution>` -
we'll run Linux applications using `chroot /compat/<distribution>`:

```
# cat /etc/jail.conf.d/jammy.conf
jammy {
  depend = ldap;

  $id = "111";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  allow.raw_sockets;
  exec.clean;

  vnet;
  mount.devfs;
  devfs_ruleset = 5;

  allow.mount;
  allow.mount.devfs;
  allow.mount.fdescfs;
  allow.mount.procfs;
  allow.mount.linprocfs;
  allow.mount.linsysfs;
  allow.mount.tmpfs;

  enforce_statfs = 1;

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";

  mount += "devfs     $path/compat/jammy/dev     devfs     rw  0 0";
  mount += "tmpfs     $path/compat/jammy/dev/shm tmpfs     rw,size=1g,mode=1777  0 0";
  mount += "fdescfs   $path/compat/jammy/dev/fd  fdescfs   rw,linrdlnk 0 0";
  mount += "linprocfs $path/compat/jammy/proc    linprocfs rw  0 0";
  mount += "linsysfs  $path/compat/jammy/sys     linsysfs  rw  0 0";
  mount += "/home     $path/compat/jammy/home    nullfs    rw  0 0";
}
```

There is no need to mount `/tmp` or `/home` unless one needs to run X11
applications.

Start the jail:

```
# service jail start jammy
```

## Verify work

```console
# jexec jammy ifconfig -a -g epair
epair4b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:jammy:bridge0
	options=8<VLAN_MTU>
	ether 02:09:e8:f9:1e:0b
	inet 192.168.1.111 netmask 0xffffff00 broadcast 192.168.1.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
epair5b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:jammy:bridge1
	options=8<VLAN_MTU>
	ether 02:6c:fc:c8:68:0b
	inet 10.0.1.111 netmask 0xffffff00 broadcast 10.0.1.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
# jexec jammy route get 0
   route to: default
destination: default
       mask: default
    gateway: unifi.localdomain
        fib: 0
  interface: epair4b
      flags: <UP,GATEWAY,DONE,STATIC>
 recvpipe  sendpipe  ssthresh  rtt,msec    mtu        weight    expire
       0         0         0         0      1500         1         0
# jexec jammy nc -uv 1.1.1.1 53
Connection to 1.1.1.1 53 port [udp/domain] succeeded!
^C
# jexec jammy chroot /compat/jammy uname -srm
Linux 5.15.0 x86_64
```


# APPLICATIONS

Linux applications reside under `/compat/jammy` in various binary folders. They
still look up for configuration files, libraries, executables, in paths with
respect to root `/`, which holds FreeBSD code.

To fix the issue, run Linux applications using chroot(8):

```console
# jexec jammy chroot /compat/jammy /bin/ls
```

## Services

Start Linux services from `/etc/rc.local`:

```console
# jexec jammy cat /etc/rc.local
chroot /compat/jammy /bin/sshd
```

The services show up in the jail's process tree with ps(1).

Keep in mind that FreeBSD service(8) does not know about Linux services:

```console
# jexec jammy \
    service sshd status
sshd is not running.
# jexec jammy \
    service syslogd status
syslogd is running as pid 97816.
# jexec jammy \
    chroot /compat/jammy /usr/sbin/service --status-all
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

## Enter jail

FreeBSD environment:

```console
# jexec jammy su -l
```

Linux environment:

```console
# jexec jammy chroot /compat/jammy su -l
```
