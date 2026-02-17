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

Use tmux(1) to attach to VM consoles
([doc](https://github.com/freebsd/vm-bhyve/wiki/Using-tmux)). Otherwise
disconnecting from consoles is non-trivial `~+Ctrl-d`.

**WARNING**: tmux(1) will run under super-user.

```console
% doas pkg install tmux
% doas vm set console=tmux
```

## FreeBSD VM

Pick up an ISO from https://download.freebsd.org/releases/ISO-IMAGES/15.0/ and
add it to vm-bhyve(8) (minimize traffic by using `.xz` archive):

```console
% doas vm iso https://download.freebsd.org/releases/ISO-IMAGES/15.0/FreeBSD-15.0-RELEASE-amd64-bootonly.iso.xz
% doas vm iso
DATASTORE           FILENAME
default             FreeBSD-15.0-RELEASE-amd64-bootonly.iso.xz
```

It didn't open the archive, extract it manually:

```console
# cd /usr/vm/.iso
# xz -d FreeBSD-15.0-RELEASE-amd64-bootonly.iso
# vm iso
DATASTORE           FILENAME
default             FreeBSD-15.0-RELEASE-amd64-bootonly.iso
```

Create a VM and install downloaded image:

```console
% doas vm create -t freebsd-zvol freebsd
```

It creates a VM with sparse ZVOL, i.e., ZFS dataset allocates data as needed
instead of pre-allocating a block as typically done with `.img` files that back
the virtual filesystem:

```console
% grep disk /usr/vm/freebsd/freebsd.conf
disk0_type="virtio-blk"
disk0_name="disk0"
disk0_dev="sparse-zvol"
```

Change RAM to 512MB or 1024MB:

```console
% doas vm edit freebsd
```

Install the OS:

```console
% doas vm install freebsd FreeBSD-15.0-RELEASE-amd64-bootonly.iso
Starting freebsd
  * found guest in /usr/vm/freebsd
  * booting...
% doas vm list
NAME     DATASTORE  LOADER     CPU  MEMORY  VNC  AUTO  STATE
freebsd  default    bhyveload  1    512M    -    No    Bootloader (994)
% doas vm console freebsd
```

FreeBSD starts up and offers to install the system or use live CD.

Notice that vm-bhyve(8) creates a separate ZFS dataset for the VM and a sub-set
to emulate the disk:

```console
% zfs list -r zroot/usr/vm/freebsd
NAME                         USED  AVAIL  REFER  MOUNTPOINT
zroot/usr/vm/freebsd         176K  1.73T   120K  /usr/vm/freebsd
zroot/usr/vm/freebsd/disk0    56K  1.73T    56K  -
```
