<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**enjailed** - enjailed BSD hypervisor

# DESCRIPTION

## Host Setup

Load Kernel modules for bhyve:

```
# load Kernel modules for bhyve
echo 'vmm_load="YES"' >> /boot/loader.conf
echo 'nmdm_load="YES"' >> /boot/loader.conf
```

Alternatively, load modules while running the OS:

```
kldload vmm
kldload nmdm
```

Create a jail to run bhyve:

```console
zfs clone zroot/usr/jail/15.0@skel zroot/usr/jail/bhyve
```

Create a new devfs(8) ruleset to access vmm(8) and nmdm(8):

```
% cat /etc/devfs.rules
[devfsrules_jail_bhyve=10]
add include $devfsrules_jail_vnet
add path vmm unhide
add path vmm/* unhide
add path tap* unhide
add path nmdm* unhide
```

Create a jail configuration with following settings:

```
# file: /etc/jail.conf.d/bhyve.conf
bhyve {
    allow.vmm;
    devfs_ruleset = 10;
}
```

## Inside Bhyve jail

Make sure the network is up and running:

```
% head /etc/resolv.conf /etc/rc.conf.d/routing
==> /etc/resolv.conf <==
nameserver 10.0.0.2

==> /etc/rc.conf.d/routing <==
defaultrouter="10.0.0.1"
```

Install vm-bhyve to manage bhyve virtual machines:

```
pkg install vm-bhyve
```

Create a local directory to store VMs:

```
mkdir /usr/vm
```

Setup vm:

```
% cat /usr/local/etc/rc.conf.d/vm
vm_enable="yes"
vm_dir="/usr/vm"
```

Initialize vm:

```
vm init
```
