<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Name

**bootstrap** - setup host to run jails

# DESCRIPTION

Store jails under `/usr/jail/` dataset.

```console
# zfs create -v zroot/usr/jail
```

Setup jail service:

```console
# sysrc -f /etc/rc.conf.d/jail jail_enable=yes
# sysrc -f /etc/rc.conf.d/jail jail_reverse_stop=yes
```

Store jail configurations in `/etc/jail.conf.d`:

```console
# mkdir /etc/jail.conf.d

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/host/etc/jail.conf
# mv -nv /tmp/jail.conf /etc/

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/host/etc/jail.conf.d/template.conf
# mv -nv /tmp/template.conf /etc/jail.conf.d/
```

## Network

Lab runs thick jails, isolated with virtual networks, one per jail. The host
creates one or more pairs of virtual Ethernet interfaces, epair(4), with a-side
kept in the host environment and b-side given to the jail. The host connects
the a-sides using bridge(4) to allow traffic between the interfaces.

Create two briddges to separate external and intra-jail traffic:

```console
# sysrc -f /etc/rc.conf.d/network ifconfig_igc0=up

# sysrc -f /etc/rc.conf.d/network cloned_interfaces='bridge0 bridge1'

# sysrc -f /etc/rc.conf.d/network ifconfig_bridge0_name="jail_ext"
# sysrc -f /etc/rc.conf.d/network ifconfig_jail_ext='inet 192.168.10.101/24 addm igc0 up descr jail:ext'

# sysrc -f /etc/rc.conf.d/network ifconfig_bridge1_name="jail_int"
# sysrc -f /etc/rc.conf.d/network ifconfig_jail_int='inet 10.0.0.1/24 up descr jail:int'
```

Use [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail)
script to manage jail epair(4) setup and tear down.

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/rc/rc.jail
# mv -nv /tmp/rc.jail /usr/local/etc/rc.jail
```

## Firewall

Allow PF to [forward packets](https://github.com/freebsd/freebsd-src/blob/be4f245e1e4fe60d43aaff5b11b45f2a9a66a51c/libexec/rc/rc.d/routing#L387-L393):

```console
# sysrc -f /etc/rc.conf.d/routing gateway_enable=yes
```

Enable PF and logging:

```console
# sysrc -f /etc/rc.conf.d/pf pf_enable=yes
# service pf start

# sysrc -f /etc/rc.conf.d/pflog pflog_enable=yes
# service pflog start
```

Add PF configuration - deny all traffic by default, pass through only what is
needed:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/host/etc/pf.conf
# mv -nv /tmp/pf.conf /etc/pf.conf

-- validate file
# pfctl -nf /etc/pf.conf
# pfctl -f /etc/pf.conf

-- show translations
# pfctl -sn
-- show filters
# pfctl -sr
```
