<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**jail** - jails provide virtual environment

# DESCRIPTION

Setup *thick* jails with full FreeBSD environment excluding Kernel. Manage jails
with builtin tools: jail(8), jls(8), and jexec(8).

- [bootstrap](./bootstrap.md): prepare the host to serve jails - create a ZFS
  dataset for jails, configure networking, setup jail(8) service.

- [template](./freebsd_template.md): create a jail template, ready for cloning
  to create new jails from in the future. It has network setup, all packages up
  to date, snapshotted.

- Ubuntu [template](./ubuntu_template.md), [jail](./ubuntu_jail.md): create a
  jail running Ubuntu using debootstrap(8). It installs Ubuntu environment under
  `/compat/<distribution>` in the jail. Run Ubuntu environment by chroot(8) at
  that folder.
