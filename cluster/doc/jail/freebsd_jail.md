NAME
====

**freebsd-jail** - create a FreeBSD jail

DESCRIPTION
===========

ZFS Bootstrap
-------------

Start from FreeBSD template:

```console
% zfs list -t snapshot zroot/jail/template/15.0-RELEASE
NAME                                    USED  AVAIL  REFER  MOUNTPOINT
zroot/jail/template/15.0-RELEASE@p0.0   104K      -   374M  -
# zfs clone zroot/jail/template/15.0-RELEASE@p0.0 zroot/jail/container/demo
```

Configuration
-------------

Use `rc.jail` script on the host environment manage jail epair(4) for every bridge:

```console
# cat /etc/jail.conf.d/demo.conf
demo {
  $id = "123";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  host.hostname = "${name}.lab.net";
  path = "/jail/container/${name}";

  vnet;
  allow.raw_sockets;

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

Verify work:

```console
% doas jexec demo ifconfig -a -g epair
epair8b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
  description: jail:demo:bridge0
  options=8<VLAN_MTU>
  ether 02:86:6e:e0:e2:0b
  inet 192.168.1.123 netmask 0xffffff00 broadcast 192.168.1.255
  groups: epair
  media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
  status: active
  nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
epair9b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
  description: jail:demo:bridge1
  options=8<VLAN_MTU>
  ether 02:7e:3f:38:56:0b
  inet 10.0.1.123 netmask 0xffffff00 broadcast 10.0.1.255
  groups: epair
  media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
  status: active
  nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
% doas jexec demo route get 0
   route to: default
destination: default
       mask: default
    gateway: 192.168.1.1
        fib: 0
  interface: epair8b
      flags: <UP,GATEWAY,DONE,STATIC>
 recvpipe  sendpipe  ssthresh  rtt,msec    mtu        weight    expire
       0         0         0         0      1500         1         0
% doas jexec demo nc -uv 192.168.1.1 53
Connection to 192.168.1.1 53 port [udp/domain] succeeded!
^C
% doas jexec demo ping -c 1 freebsd.org
PING freebsd.org (96.47.72.84): 56 data bytes
64 bytes from 96.47.72.84: icmp_seq=0 ttl=51 time=32.784 ms

--- freebsd.org ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 32.784/32.784/32.784/0.000 ms
```

Start jail
----------

Manually:

```console
# service jail start demo
```

Auto-start:

```console
# sysrc -f /etc/rc.conf.d/jail jail_list+=demo
jail_list: '' -> demo
```

SEE ALSO
========

-	https://freebsdfoundation.org/wp-content/uploads/2020/03/Jail-vnet-by-Examples.pdf
