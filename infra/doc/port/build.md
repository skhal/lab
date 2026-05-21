<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**build** -- build packages using poudriere(8)

# DESCRIPTION

## Create a jail

Run a jail:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/port/data/host/etc/jail.conf.d/pkg.conf
# mv -nv /tmp/pkg.conf /etc/jail.conf.d/

# service jail run pkg
```

## Prepare environment

Install development version of poudriere (it comes with a number of fixes
including poudriere run in a jail). Keep in mind, that we add other
dependencies because the jail template is minimal:

```console
# pkg install poudriere-devel git FreeBSD-jail FreeBSD-clibs-dev FreeBSD-mtree FreeBSD-bmake FreeBSD-clang ccache

# grep DISTFILES_CACHE /usr/local/etc/poudriere.conf
DISTFILES_CACHE=/usr/ports/distfiles
# mkdir -p /usr/ports/distfiles

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/port/data/jail/usr/local/etc/poudriere.conf.diff
# patch -lb -i /tmp/poudriere.conf.diff /usr/local/etc/poudriere.conf
```

Initialize poudriere(8) and create a build jail:

```console
# poudriere ports -c

# poudriere jail -c -j 15amd64 -v 15.0-RELEASE
```

## Build packages

Use a list of packages to build:

```console
# cat /usr/local/etc/poudriere.d/pkglist
devel/protobuf

# poudriere bulk -j 15amd64 -J 16 -f /usr/local/etc/poudriere.d/pkglist
```
