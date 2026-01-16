<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**jail** - setup jails

Summary
-------

This section describes how to setup and manage thick jails using basic tools, jail(8) and jexec(8).

There are two parts to jails setup:

1.	[Bootstrap](./bootstrap.md): prepare the OS to host jails. The steps include create ZFS datasets, configure the networking, setup the jail service.

2.	Create and manage jails. The jails split into a template and a running container using ZFS datasets. Such setup allows quick spawn of jails by cloning templates into new containers.

Start with FreeBSD jail instructions to create an isolated user lands using freebsd-base(7):

-	[FreeBSD template](./freebsd_template.md)
-	[FreeBSD jail](./freebsd_jail.md)

A Linux environment departures from the FreeBSD setup. It adds Linux Compatibility layer to the FreeBSD setup, bootstrapped at `/compat/<distribution>` prefix, to provide Linux user land:

-	[Ubuntu template](./ubuntu_template.md)
-	[Ubuntu jail](./ubuntu_jail.md)
