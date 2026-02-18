<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**enjailed** - enjailed BSD hypervisor

# DESCRIPTION

## Jail bootstrap

Create an empty VNET jail:

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

## Enajil ZFS

The goal is to store VMs in a ZFS dataset under jail management.

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

Enable ZFS service inside the jail to mount the dataset:

```console
# pkg -r /usr/jail/bhyve/ install FreeBSD-zfs
# sysrc -f /usr/jail/bhyve/etc/rc.conf.d/zfs zfs_enable=yes
zfs_enable: NO -> yes
```

jail(8) can
[attach](https://github.com/freebsd/freebsd-src/blob/349808d8bd197165390a286bccdaa29a1d77c7ab/usr.sbin/jail/jail.c#L101)
ZFS datasets during jail startup, right after `exec.created` runs. It picks up
datasets from `zfs.dataset` param.

Note that the dataset to be enjailed must have `jailed` property set:
[check](https://github.com/freebsd/freebsd-src/blob/349808d8bd197165390a286bccdaa29a1d77c7ab/usr.sbin/jail/command.c#L614)
.

There is also no symmetric zfs-unjail(8) option. There is no way to run the
command manually in easy way because it needs to be run on the host while jail
is still available. Stop commands
[sequence](https://github.com/freebsd/freebsd-src/blob/e1e18cc12e68762b641646b203d9ac42d10e3b1f/usr.sbin/jail/jail.c#L112-L114)
does not have a hook for that: it runs `exec.prestop` in the host environment,
then `exec.stop` in the jail environment, and executes `jail stop`. Commands
to be run after `jail stop` don't have access to the jail because it was
removed, resulting in "invalid jail name".

```console
% grep zfs /etc/jail.conf.d/bhyve.conf
  allow.mount.zfs;
  zfs.dataset = "zroot/usr/jail/${name}/vm";
```

Include ZFS service in the shutdown sequence to unmount the VMs dataset:

```console
# diff -u /usr/jail/bhyve/etc/rc.d/zfs{.orig,}
--- /usr/jail/bhyve/etc/rc.d/zfs.orig	2026-02-18 14:10:06.114400000 -0600
+++ /usr/jail/bhyve/etc/rc.d/zfs	2026-02-18 14:10:24.899516000 -0600
@@ -5,6 +5,7 @@
 # PROVIDE: zfs
 # REQUIRE: zfsbe
 # BEFORE: FILESYSTEMS var
+# KEYWORD: shutdown

 . /etc/rc.subr
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
