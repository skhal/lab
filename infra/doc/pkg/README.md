<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**pkg** -- pkg(8) repository of custom ports

# DESCRIPTION

## Create a jail

Run a jail:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/pkg/data/host/etc/jail.conf.d/pkg.conf
# mv -nv /tmp/pkg.conf /etc/jail.conf.d/

# service jail run pkg
```

## Prepare environment

Install development version of poudriere (it comes with a number of fixes
including poudriere run in a jail). Keep in mind, that we add other
dependencies because the jail template is minimal:

```console
# pkg install poudriere-devel \
    FreeBSD-bmake \
    FreeBSD-clang \
    FreeBSD-clibs-dev \
    FreeBSD-jail \
    FreeBSD-mtree \
    ccache \
    git

# grep DISTFILES_CACHE /usr/local/etc/poudriere.conf
DISTFILES_CACHE=/usr/ports/distfiles
# mkdir -p /usr/ports/distfiles

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/pkg/data/jail/usr/local/etc/poudriere.conf.diff
# patch -lb -i /tmp/poudriere.conf.diff /usr/local/etc/poudriere.conf
```

Initialize poudriere(8) and create a build jail:

```console
# poudriere ports -c

# poudriere jail -c -j 15amd64 -v 15.0-RELEASE
```

Generate a key to sign packages -- it will be used to connect to the server:

```console
# mkdir /usr/local/etc/poudriere.d/keys

# openssl genrsa -out /usr/local/etc/poudriere.d/keys/pkg.key 4096
# chmod 400 /usr/local/etc/poudriere.d/keys/pkg.key
# openssl rsa -in /usr/local/etc/poudriere.d/keys/pkg.key -pubout -out /usr/local/etc/poudriere.d/keys/pkg.pub
```

Point `PKG_REPO_SIGNING_KEY` at pkg.key in `poudriere.conf`. Now on,
poudriere(8) signs repo during builds.

## Build packages

Use a list of packages to build:

```console
# cat /usr/local/etc/poudriere.d/pkglist
devel/protobuf

# poudriere bulk -j 15amd64 -J 16 -f /usr/local/etc/poudriere.d/pkglist
```

## Upgrade a package

Example: want to update devel/protobuf package from 29.6 to 34.2.

The package uses cmake to build and requires additional tools or libraries.

```console
# pkg install \
    FreeBSD-clang-dev \
    FreeBSD-toolchain \
    cmake
```

Work in the default ports tree:

```console
# cd /usr/local/poudriere/ports/default/
# git checkout --track -b protobuf-v34.2

# cd devel/protobuf
```

Update:

- Makefile: bump the version
- distinfo: change the release archive, sha256, size, and set timestamp to
  `date -u +%s`

```console
# make
```

If the package builds successfully, re-generate package list (plist):

```console
# makeplist
```

Review the generated file and remove the first todo line that says: `/you/...`.
The package is ready to be build by poudriere(8):

```console
# poudriere bulk -j 15amd64 -J 16 -f /usr/local/etc/poudriere.d/pkglist
```

## Server

```console
# pkg install nginx

# mkdir /usr/local/etc/rc.conf.d
# sysrc -f /usr/local/etc/rc.conf.d/nginx ngingx_enable=yes

# service nginx start
```

Setup the configuration to serve the clients by `${ABI}` which is equivalent to
`FreeBSD:15:amd64` or `${OS}:${OS_VERSION}:${ARCH}`:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/pkg/data/jail/usr/local/etc/nginx/nginx.conf.diff
# patch -lb -i /tmp/nginx.conf.diff /usr/local/etc/nginx/nginx.conf

# service nginx reload
```

## Client

Copy public key from the server to the client machine:

```console
# mkdir /usr/local/etc/pkg/repos/lab/
# cp /tmp/pkg.pub /usr/local/etc/pkg/repos/lab/

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/pkg/data/client/usr/local/etc/pkg/repos/lab.conf
# mv -nv /tmp/lab.conf /usr/local/etc/pkg/repos/
```

Fetch the repository:

```console
# pkg update -r lab
# pkg search -r lab protobuf
# pkg install -r lab protobuf
```
