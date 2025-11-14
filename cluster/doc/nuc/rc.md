Name
====

**rc** - FreeBSD resource configuration

Description
===========

FreeBSD boot uses init(8) to trigger rc(8) that in turn start services in topological order (see rcorder(8)).

Services
--------

Services use a management scripts to configure dependencies, rc-variables (aka flags), actions (start, stop, etc.) and functions to execute actions.

The scripts reside in `rc.d` folders. RC is flexible to support multiple locations:

-	Default: `/etc/rc.d`, is for services that come with FreeBSD installation.
-	Local: `/usr/local/etc/rc.d`, is for services installed with pkg(1).

There is a flag to let RC know the location of local RC scripts:

```console
% sysrc local_startup
local_startup: /usr/local/etc/rc.d
```

Dependency graph
----------------

Get a list of services available in topological order:

```console
# service -r
```

Or only enabled services in topological order:

```console
# service -e
```

Sometimes it is helpful to analyse the dependency graph by for standard services only or group the services that can start in parallel with `-p`. Use rcorder(8):

```console
% rcorder -p /etc/rc.d/*
/etc/rc.d/dhclient /etc/rc.d/dumpon /etc/rc.d/dnctl /etc/rc.d/natd /etc/rc.d/sysctl
/etc/rc.d/ddb /etc/rc.d/hostid
...
```

Flags
-----

RC flags can be found in rc.conf(5) files. The flags configure services. RC loads configuration files in the following order:

-	Default values from `/etc/default/rc.conf`.
-	Vendor overrides from `/etc/default/vendor.conf`.
-	Global overrides `/etc/rc.conf` and `/etc/rc.conf.local` (legacy).

There are per-service overrides that apply on top of global flags:

-	Service overrides `<dir>/rc.conf.d/<name>` for the service `<name>`, where`<dir>/` is one of the `rc.d/` prefixes, e.g. `/etc` or `/usr/local`.

Note that the last loaded flag value wins.

Configuration files
-------------------

The `/etc/rc.conf` file is a global container of flags, shared between all services and rc(8) itself.

It is a good practice to limit the visibility of flags by placing the flags into as narrow as possible scope. That is place service flags under `/etc/rc.conf.d/<name>`.

Use sysrc(8) to get a list of files:

```console
% sysrc -s hostname -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/hostname /usr/local/etc/rc.conf.d/hostname
```

*Warning*: some services share settings, i.e., `dhclient` and `netif` share DHCP settings, to be stored in `/etc/rc.conf.d/network`:

```console
% sysrc -s dhclient -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/dhclient /usr/local/etc/rc.conf.d/dhclient /etc/rc.conf.d/network /usr/local/etc/rc.conf.d/network
% sysrc -s netif -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/network /usr/local/etc/rc.conf.d/network /etc/rc.conf.d/netif /usr/local/etc/rc.conf.d/netif
```
