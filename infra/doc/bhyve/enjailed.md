<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**enjailed** - enjailed BSD hypervisor

# DESCRIPTION

## Jail setup

First, create an empty jail with VNET:

```console
# zfs clone zroot/usr/jail/15.0@skel zroot/usr/jail/bhyve
```

Configuration:

```console
% # cat /etc/jail.conf.d/bhyve.conf
% doas service jail start bhyve
bhyve {
  $id = "6";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.10.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.0.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  host.hostname = "${name}.lab.net";
  path = "/usr/jail/${name}";

  depend  = "dns";
  depend += "ldap";

  # keep-sorted start
  allow.raw_sockets;
  devfs_ruleset = 5;
  enforce_statfs = 1;
  mount.devfs;
  vnet;
  # keep-sorted end

  exec.clean;
  exec.consolelog = "/var/log/jail_${name}.log";

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.start    = "/bin/sh /etc/rc";

  exec.stop     = "/bin/sh /etc/rc.shutdown";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```

Verify network connection:

```console
% doas jexec bhyve su -
bhyve # route get 0
   route to: default
destination: default
       mask: default
    gateway: 10.0.0.1
        fib: 0
  interface: epair8b
      flags: <UP,GATEWAY,DONE,STATIC>
 recvpipe  sendpipe  ssthresh  rtt,msec    mtu        weight    expire
       0         0         0         0      1500         1         0
bhyve # drill -Q freebsd.org
96.47.72.84
bhyve # ping -c 1 freebsd.org
PING freebsd.org (96.47.72.84): 56 data bytes
64 bytes from 96.47.72.84: icmp_seq=0 ttl=50 time=29.651 ms

--- freebsd.org ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 29.651/29.651/29.651/0.000 ms
```

Next, add a ZFS data for VMs and delegate it to the jail for management. The
goal is to let bhyve jail manage the dataset: mount, create sub-datasets, etc.

According to zfs-jail(8), a dataset must meet following requirements before a
jail can manage it:

1. `jailed` property must be set.
2. devfs.conf(5) entries must expose `/dev/zfs` device within the jail.
3. Jail configuration must have `allow.mount` and `allow.mount.zfs` parameters
   set.
4. Jail must nave `enforce_statfs` set to a value lower than 2.

Bhyve jail configuration meets the `statfs`-requirements:

```console
% grep enforce_statfs /etc/jail.conf.d/bhyve.conf
  enforce_statfs = 1;
```

The `devfs`-item is also met:

```console
% grep devfs_ruleset /etc/jail.conf.d/bhyve.conf
  devfs_ruleset = 5;
```

with the following ruleset (redacted to keep only relevant items):

```
[devfsrules_jail=4]
add path zfs unhide

[devfsrules_jail_vnet=5]
add include $devfsrules_jail
```

Add `allow.mount` and `allow.mount.zfs` to jail configuration:

```console
% grep allow.mount /etc/jail.conf.d/bhyve.conf
  allow.mount.zfs;
  allow.mount;
```

Create a ZFS dataset for VMs:

```console
# zfs create -o jailed=on zroot/usr/jail/bhyve/vm
```

Verify:

```console
# zfs get jailed,canmount,mounted,mountpoint -r zroot/usr/jail/bhyve
NAME                     PROPERTY    VALUE               SOURCE
zroot/usr/jail/bhyve     jailed      off                 default
zroot/usr/jail/bhyve     canmount    on                  default
zroot/usr/jail/bhyve     mounted     yes                 -
zroot/usr/jail/bhyve     mountpoint  /usr/jail/bhyve     inherited from zroot/usr
zroot/usr/jail/bhyve/vm  jailed      on                  local
zroot/usr/jail/bhyve/vm  canmount    on                  default
zroot/usr/jail/bhyve/vm  mounted     no                  -
zroot/usr/jail/bhyve/vm  mountpoint  /usr/jail/bhyve/vm  inherited from zroot/usr
```

Before we test the setup within the jail, we need to add ZFS support to the
jail and enable ZFS service to manage the mounts:

```console
# pkg -r /usr/jail/bhyve/ install FreeBSD-zfs
# sysrc -f /usr/jail/bhyve/etc/rc.conf.d/zfs zfs_enable=yes
zfs_enable: NO -> yes
```

At this point everything is set up except the fact that ZFS dataset is not
enjailed with zfs-jail(8).

Apparently jail(8) can
[attach ZFS datasets](https://github.com/freebsd/freebsd-src/blob/349808d8bd197165390a286bccdaa29a1d77c7ab/usr.sbin/jail/jail.c#L101)
during the start command sequence, right after `exec.created` runs, if jail
configuration has `zfs.dataset` param set. It
[checks](https://github.com/freebsd/freebsd-src/blob/349808d8bd197165390a286bccdaa29a1d77c7ab/usr.sbin/jail/command.c#L614)
that the ZFS dataset has `jailed` property set.

Keep in mind that there is no analogous step to zfs-unjail(8) the dataset. We
need to manually run it. Based on where zfs-jail(8) runs in the start sequence,
natural place in the stop sequence is `exec.poststop`:

```console
% grep zfs /etc/jail.conf.d/bhyve.conf
  allow.mount.zfs;
  zfs.dataset = "zroot/usr/jail/${name}/vm";
  exec.poststop += "/sbin/zfs unjail zroot/usr/jail/${name}/vm";
```

Verify work:

```console
% doas service jail start bhyve
% doas jexec bhyve ls /dev/zfs
/dev/zfs
% doas jexec bhyve zfs mount
zroot/usr/jail/bhyve            /
zroot/usr/jail/bhyve/vm         /usr/jail/bhyve/vm
```

At last, lets change the VMs mount point inside the jail (keep in mind that ZFS auto-mounts the dataset when `mountpoint` property is set):

```console
% doas jexec bhyve zfs umount zroot/usr/jail/bhyve/vm
% doas jexec bhyve zfs set mountpoint=/vm zroot/usr/jail/bhyve/vm
% doas jexec bhyve zfs mount
zroot/usr/jail/bhyve            /
zroot/usr/jail/bhyve/vm         /vm
% doas jexec bhyve rm -rf /usr/jail/
```
