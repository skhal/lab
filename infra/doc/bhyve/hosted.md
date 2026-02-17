<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**hosted** - hosted BSD hypervisor

# DESCRIPTION

This page describes how to run FreeBSD hypervisor bhyve(8) on the host (hosted
environment).

## Host setup

*Note*: create a new boot environment bectl(8) to experiment with bhyve(8).

Use vm-bhyve(8) (https://github.com/freebsd/vm-bhyve) instead of low-level
commands bhyve(8), bhyvectl(8), and bhyveload(8) to manage bhyve(8):

```console
% doas pkg install vm-bhyve
```

*Note*: `/usr/local/bin/vm` is a shell script. It relies on some of the rc(8)
configurations. Enable debug level as needed:

```console
% doas sysrc rc_debug="yes"
```

Configure doas(1):

```console
# echo 'permit nopass setenv { EDITOR=nvim } op cmd vm' >> /usr/local/etc/doas.conf
```

Isolate bhyve template, OS images, and running virtual machines (VMs) in a ZFS
dataset:

```console
# zfs create zroot/usr/vm
```

Enable vm-bhyve(8) and tell it the location of VMs:

```console
% doas sysrc -f /usr/local/etc/rc.conf.d/vm vm_enable=yes
% doas sysrc -f /usr/local/etc/rc.conf.d/vm vm_dir="zfs:zroot/usr/vm"
```

Initialize vm-bhyve(8):

```console
% doas vm init
```

Pick up example template configurations:

```console
# cp /usr/local/share/examples/vm-bhyve/* /usr/vm/.templates/
```

vm-bhyve(8) uses [virtual switches](https://github.com/freebsd/vm-bhyve/wiki/Virtual-Switches)
for networking. A typical setup creates a virtual switch and adds a network
interface (IF) with outbound connection.

Lab uses if_bridge(8) to connect jails to internet. It means that outbound IF
is already part of the bridge and can't be added to the virtual switch.

Create a virtual switch with bridge interface and call it `public` (VM templates
use `public` switch for networking):

```console
% doas vm switch create -t manual -b brdige0 public
```

Verify:

```console
% doas vm switch list
NAME    TYPE    IFACE    ADDRESS  PRIVATE  MTU  VLAN  PORTS
public  manual  bridge0  n/a      no       n/a  n/a   n/a
% doas vm switch info public
------------------------
Virtual Switch: public
------------------------
  type: manual
  ident: bridge0
  vlan: -
  physical-ports: -
  bytes-in: 0 (OB)
  bytes-out: 0 (OB)
```
