<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**alpine-linux** - enjail [Alpine Linux](https://www.alpinelinux.org)

# DESCRIPTION

```console
# zfs clone zroot/usr/jail/15.0@skel zroot/usr/jail/alpine
# cat /etc/jail.conf.d/alpine.conf
alpine {
  $id = "13";

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
  devfs_ruleset = 5;
  enforce_statfs = 1;
  mount.devfs;
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

Verify the jail works and has network access:

```console
% doas service jail start alpine
% jls -j alpine
   JID  IP Address      Hostname                      Path
     7                  alpine.lab.net                /usr/jail/alpine
% doas jexec alpine netstat -4r
Routing tables

Internet:
Destination        Gateway            Flags         Netif Expire
default            10.0.0.1           UGS         epair7b
10.0.0.0/24        link#33            U           epair7b
10.0.0.13          link#34            UHS             lo0
localhost          link#34            UH              lo0
% doas jexec alpine drill -Q alpinelinux.org
213.219.36.190
% doas jexec alpine ping -c 1 alpinelinux.org
PING alpinelinux.org (213.219.36.190): 56 data bytes
64 bytes from 213.219.36.190: icmp_seq=0 ttl=44 time=107.647 ms

--- alpinelinux.org ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 107.647/107.647/107.647/0.000 ms
```

Install Alpine environment:

```console
% doas jexec alpine su -
root@alpine:~ # mkdir -vp /compat/alpine
/compat
/compat/alpine
root@alpine:~ # fetch -o /tmp https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/x86_64/alpine-minirootfs-3.23.3-x86_64.tar.gz
/tmp/alpine-minirootfs-3.23.3-x86_64.tar.gz           3626 kB 9921 kBps    00s
root@alpine:~ # tar -C /compat/alpine -xzf /tmp/alpine-minirootfs-3.23.3-x86_64.tar.gz
```

If we try to run any of the alpine binaries, they'll fail with high probability
because the jail does not mount a number of virtual filesystems for Linux
compatibility layer.

Shutdown the jail and make following changes:

```console
% diff -u /etc/jail.conf.d/alpine.conf{.orig,}
--- /etc/jail.conf.d/alpine.conf.orig	2026-02-20 14:15:13.011356000 -0600
+++ /etc/jail.conf.d/alpine.conf	2026-02-20 14:16:58.289166000 -0600
@@ -29,6 +29,14 @@
   mount.devfs;
   # keep-sorted end

+  # keep-sorted start
+  mount += "devfs     $path/compat/alpine/dev     devfs     rw 0 0";
+  mount += "fdescfs   $path/compat/alpine/dev/fd  fdescfs   rw,linrdlnk 0 0";
+  mount += "linprocfs $path/compat/alpine/proc    linprocfs rw 0 0";
+  mount += "linsysfs  $path/compat/alpine/sys     linsysfs  rw 0 0";
+  mount += "tmpfs     $path/compat/alpine/dev/shm tmpfs     rw,size=1g,mode=1777 0 0";
+  # keep-sorted end
+
   exec.clean;
   exec.consolelog = "/var/log/jail_${name}.log";
```

Verify it works:

```console
% doas jexec alpine chroot /compat/alpine su -
alpine:~# echo 'nameserver 10.0.0.2' > /etc/resolv.conf
alpine:~# ping -c 1 alpinelinux.org
PING alpinelinux.org (213.219.36.190): 56 data bytes
64 bytes from 213.219.36.190: seq=0 ttl=44 time=106.767 ms

--- alpinelinux.org ping statistics ---
1 packets transmitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 106.767/106.767/106.767 ms
```
