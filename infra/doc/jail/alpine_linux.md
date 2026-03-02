<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**alpine-linux** - enjail [Alpine Linux](https://www.alpinelinux.org)

# DESCRIPTION

**TL;DR** create a thin jail with Alpine Linux jail without FreeBSD jailed
environment.

```console
# zfs create zroot/usr/jail/alpine
# fetch -o /tmp https://dl-cdn.alpinelinux.org/alpine/latest-stable/releases/x86_64/alpine-minirootfs-3.23.3-x86_64.tar.gz
# tar -C /usr/jail/alpine -xzf /tmp/alpine-minirootfs-3.23.3-x86_64.tar.gz
```

Configure:

```console
# echo 'nameserver 10.0.0.2' >/usr/jail/alpine/etc/resolv.conf
# echo "auto lo" > /usr/jail/alpine/etc/network/interfaces
# chroot /usr/jail/alpine env PS1='alpine # ' /bin/sh
alpine # apk update
alpine # apk upgrade
alpine # apk add openrc
alpine # mkdir /run/openrc
alpine # touch /run/openrc/softlevel
```

Initialize passwords:

```console
# cd /usr/jail/alpine/etc/
# echo 'root::0:0::0:0:Charlie &:/root:/bin/sh' > master.passwd
# pwd_mkdb -d . -p master.passwd
```

Create a jail configuration:

```console
% cat /etc/jail.conf.d/alpine.conf
alpine {
  $id = "13";

  $gateway = "10.0.0.1";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.0.${id}/24";

  $bridges = "${bridge1}";
  $bridgeips = "${bridge1_ip}";

  host.hostname = "${name}.lab.net";
  path = "/usr/jail/${name}";

  depend  = "dns";
  depend += "ldap";

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

  mount += "devfs     $path/dev     devfs     rw 0 0";
  mount += "fdescfs   $path/dev/fd  fdescfs   rw,linrdlnk 0 0";
  mount += "linprocfs $path/proc    linprocfs rw 0 0";
  mount += "linsysfs  $path/sys     linsysfs  rw 0 0";
  mount += "tmpfs     $path/dev/shm tmpfs     rw,size=1g,mode=1777 0 0";
  mount.devfs;
  devfs_ruleset = 5;

  enforce_statfs = 1;

  exec.clean;
  exec.consolelog = "/var/log/jail_${name}.log";

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.created += "/sbin/route -j ${name} add default ${gateway}";
  exec.start    = "/sbin/openrc";

  persist;

  exec.stop     = "/sbin/openrc shutdown";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```

Start the jail and verify it works:

```console
% doas ifconfig -j alpine -ag epair
epair8b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:alpine:bridge1
	options=60000b<RXCSUM,TXCSUM,VLAN_MTU,RXCSUM_IPV6,TXCSUM_IPV6>
	ether 58:9c:fc:10:f6:4a
	inet 10.0.0.13 netmask 0xffffff00 broadcast 10.0.0.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
op@nuc:~ % doas route -j alpine get 0
   route to: default
destination: default
       mask: default
    gateway: 10.0.0.1
        fib: 0
  interface: epair8b
      flags: <UP,GATEWAY,DONE,STATIC>
 recvpipe  sendpipe  ssthresh  rtt,msec    mtu        weight    expire
       0         0         0         0      1500         1         0
% doas jexec alpine ping -c 1 alpinelinux.org
PING alpinelinux.org (213.219.36.190): 56 data bytes
64 bytes from 213.219.36.190: seq=0 ttl=44 time=108.967 ms

--- alpinelinux.org ping statistics ---
1 packets transmitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 108.967/108.967/108.967 ms
```

## Install Bazel

```console
# apk update
# apk upgrade
# echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
# apk update
# apk add bazel8
```
