<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**jail** - isolated virtual environment

# DESCRIPTION

FreeBSD jails combine chroot(8) with resource confinment (CPU, memory, network).

A jail may run a single process (*thin*) or entire world (*thick*). This section
focuses on thick jails.

A thick jail, provides virtual envirotnment to run applications. It still
shares Kernel with the host but has own world and resource isolation - CPU,
Memory, optional network stack (Virtual Network - VNET). It uses resource
configuration (RC) scripts to create runtime environment.

Instructions:

- [bootstrap](./bootstrap.md): setup a FreeBSD host to run jails.

- [FreeBSD](./freebsd.md): a thick jail to run FreeBSD environment.

- Ubuntu [template](./ubuntu_template.md), [jail](./ubuntu_jail.md): create a
  jail running Ubuntu using debootstrap(8). It installs Ubuntu environment under
  `/compat/<distribution>` in the jail. Run Ubuntu environment by chroot(8) at
  that folder.

- [Alpine Linux](./alpine_linux.md): enjail Alpine Linux

# SEE ALSO

jail(8) • jls(8) • jexec(8)
