# NAME

**jail-ubuntu** - create a Ubuntu jail


# DESCRIPTION

We'll create a Ubuntu jail with Virtual Network and dependency on `ldap` jail
using Ubuntu template.

## ZFS Bootstrap

Start from Ubuntu template

```console
# zfs clone zroot/jail/template/Ubuntu-22.04@p2.0 zroot/jail/container/jammy
```

## Jail configuration

Propagate mount points, created by `linux` service, from the host to jailed
environment using `mount` parameter. Ensure to set mount points under
`/compat/<distribution>`.

```
# cat /etc/jail.conf.d/jammy.conf
jammy {
  $id = "111";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  host.hostname = "${name}.lab.net";
  path = "/jail/container/${name}";

  depend = ldap;

  vnet;

  # keep-sorted start
  allow.mount.devfs;
  allow.mount.fdescfs;
  allow.mount.linprocfs;
  allow.mount.linsysfs;
  allow.mount.procfs;
  allow.mount.tmpfs;
  allow.mount;
  allow.raw_sockets;
  # keep-sorted end

  mount += "/home     $path/compat/jammy/home    nullfs    rw 0 0";
  mount += "devfs     $path/compat/jammy/dev     devfs     rw 0 0";
  mount += "fdescfs   $path/compat/jammy/dev/fd  fdescfs   rw,linrdlnk 0 0";
  mount += "linprocfs $path/compat/jammy/proc    linprocfs rw 0 0";
  mount += "linsysfs  $path/compat/jammy/sys     linsysfs  rw 0 0";
  mount += "tmpfs     $path/compat/jammy/dev/shm tmpfs     rw,size=1g,mode=1777 0 0";
  mount.devfs;
  devfs_ruleset = 5;

  enforce_statfs = 1;

  exec.clean;
  exec.consolelog = "/var/log/jail_${name}.log";

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.start    = "/bin/sh /etc/rc";

  exec.stop     = "/bin/sh /etc/rc.shutdown";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```

> [!IMPORTANT]
> Mount `/tmp` and `/home` to run X11 applications.

## Start Jail

Start the jail:

```
# service jail start jammy
```

Verify work:

```console
# jexec jammy ifconfig -a -g epair
epair4b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:jammy:bridge0
	options=8<VLAN_MTU>
	ether 02:f5:a7:8f:a1:0b
	inet 192.168.1.111 netmask 0xffffff00 broadcast 192.168.1.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
epair5b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:jammy:bridge1
	options=8<VLAN_MTU>
	ether 02:86:6e:e0:e2:0b
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
# jexec jammy ping -c 1 192.168.1.1
PING 192.168.1.1 (192.168.1.1): 56 data bytes
64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=0.601 ms

--- 192.168.1.1 ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 0.601/0.601/0.601/0.000 ms
# jexec jammy nc -uv 1.1.1.1 53
Connection to 1.1.1.1 53 port [udp/domain] succeeded!
^C
# jexec jammy freebsd-version
14.3-RELEASE-p2
# jexec jammy chroot /compat/jammy uname -srm
Linux 5.15.0 x86_64
```

## Linux applications

The jail is pre-set with Linux shared libraries under `/compat/jammy`, including
Linux file system:

```console
# jexec jammy ls /compat/jammy
bin	etc	lib32	media	proc	sbin	tmp
boot	home	lib64	mnt	root	srv	usr
dev	lib	libx32	opt	run	sys	var
```

Linux applications expect to find supporting configurations, libraries at the
root file system:

```console
# jexec jammy ldd /compat/jammy/bin/ls
/compat/jammy/bin/ls:
	libselinux.so.1 => not found (0)
	libc.so.6 => not found (0)
	[vdso] (0x21a429674000)
# jexec jammy /compat/jammy/bin/ls
ELF interpreter /lib64/ld-linux-x86-64.so.2 not found, error 2
Abort
# jexec jammy ls /compat/jammy/lib64/ld-linux-x86-64.so.2
/compat/jammy/lib64/ld-linux-x86-64.so.2
```

Use chroot(8) to run Linux applications:

```console
# jexec jammy chroot /compat/jammy /bin/ls
bin   dev  home  lib32	libx32	mnt  proc  run	 srv  tmp  var
boot  etc  lib	 lib64	media	opt  root  sbin  sys  usr
```

## Linux Services

Jail starts with `/bin/sh /etc/rc` that is not aware of Linux applications or
services.

FreeBSD RC includes a
[`local`](https://github.com/freebsd/freebsd-src/blob/5f1f7d8457d4fc28c6cff7e26a629a2d6ee3fc61/libexec/rc/rc.d/local)
script to run local customizations, placed into `/etc/rc.local` and
`/etc/rc.shutdown.local`:

```console
# jexec jammy rcorder /etc/rc.d/* | cat -n | grep -C 1 '/local$'
   148	/etc/rc.d/lpd
   149	/etc/rc.d/local
   150	/etc/rc.d/hcsecd
```

For example, manage Linux SSH server:

```console
# jexec jammy cat /etc/rc.local
/usr/sbin/chroot /compat/jammy /bin/sshd
```

The service will show up in the jail process tree. FreeBSD service(8) won't
see Linux services but Linux service(8) will:

```console
# jexec jammy \
    service sshd status
sshd is not running.
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


# SEE ALSO

* https://github.com/freebsd/freebsd-src/blob/5f1f7d8457d4fc28c6cff7e26a629a2d6ee3fc61/libexec/rc/rc.d/local
