<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**jailed** - run BSD hypervisor (bhyve) in a FreeBSD jail

# DESCRIPTION

## On the host

Add VMM and NMDM kernel modules to the system start:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/bhyve/data/boot/loader.conf.diff
# patch -lb -i /tmp/loader.conf.diff /boot/loader.conf
```

Load the modules at runtime without restart:

```console
# kldload vmm
# kldload nmdm
```

Unhide VMM inside the jail:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/bhyve/data/etc/devfs.rules
# cp -n /tmp/devfs.rules /etc/devfs.rules
```

Create a jail:

```console
# zfs snapshot zroot/usr/jail/15.0@`date +%y%m%d`
# zfs clone zroot/usr/jail/15.0@260513 zroot/usr/jail/bhyve
# cp /etc/jail.conf.d/bsd.conf /etc/jail.conf.d/bhyve.conf
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/bhyve/data/etc/jail.conf.d/bhyve.conf.diff
# patch -lb -i /tmp/bhyve.conf.diff /etc/jail.conf.d/bhyve.conf
```

## In the jail

Store VMs under `/usr/bhyve`:

```console
# mkdir /usr/bhyve
```

Install bhyve software:

- use tmux(1) to connect to VM console
- install FreeBSD-bhyve and FreeBSD-acpi if the jail is minimal

```console
# pkg install -y vm-bhyve FreeBSD-bhyve FreeBSD-acpi tmux

# mkdir /usr/local/etc/rc.conf.d
# sysrc -f /usr/local/etc/rc.conf.d/vm vm_enable=yes
# sysrc -f /usr/local/etc/rc.conf.d/vm vm_dir=/usr/bhyve

# vm init
# vm set console=tmux
```

Linux VMs may start using GRUB or UEFI. The latter one uses
`/usr/local/share/uefi-firmware/BHYVE_UEFI.fd`:

```console
# pkg install -y grub2-bhyve bhyve-firmware
```

Create local network 10.0.1.0/24 to connect the VMs using if_bridge(4).
Include jail's external interface, the b-side of the epair to let the VMs
connect to the outside of the jail:

```console
# ifconfig -g epair
epair7b
# sysrc -f /etc/rc.conf.d/network cloned_interfaces="bridge0"
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge0="inet 10.0.1.1/24 addm epair7b up descr vm:ext"
```

Use packer filter (PF) to forward packets from the local network to the epair:

```console
# sysrc -f /etc/rc.conf.d/routing gateway_enable=yes

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/bhyve/data/etc/pf.conf

# sysrc -f /etc/rc.conf.d/pf pf_enable=yes
```

Use newly create bridge in bhyve:

```console
# vm switch create -t manual -b bridge0 public
```

Restart the jail for the network service to take effect. An alternative is to
restart networking and routing services.
