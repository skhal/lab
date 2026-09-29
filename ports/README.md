<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**ports**: FreeBSD ports overlay

# DESCRIPTION

This collection ports Lab software to FreeBSD and defines additional ports
that are not available in FreeBSD-ports tree.

Use poudriere(8) to build default FreeBSD-ports with Lab ports act as an
*overlay* layer, i.e. extend default ports tree.

## Build a port with poudriere(8)

Install `poudriere-devel` package (`poudriere` does not work inside a minimal
jail as of Sep '26):

- the configuration is in `/usr/local/etc/poudriere`.
- poudriere(8) builds ports tree in jails. It keeps all data under
  `/usr/local/poudriere`.

Foollow standard instructions to create a jail and checkout default FreeBSD
ports.

Clone Lab repository:

```
# cd /usr/local/poudriere/ports

# git clone https://github.com/skhal/lab
```

Add Lab ports tree to poudriere and call it `lab`:

```
# poudriere ports -c -M /usr/local/poudriere/ports/lab/ports -p lab -m null

# poudriere ports -l
PORTSTREE METHOD    TIMESTAMP           PATH
default   git+https 2026-09-13 17:00:37 /usr/local/poudriere/ports/default
lab       null      2026-09-25 14:51:01 /usr/local/poudriere/ports/lab/ports
```

Build one of the Lab ports using ports overlay:

```
# poudriere bulk -j 15amd64 -O lab -C sysutils/certval
```

## Develop a port with poudriere(8)

Create a new port under `/usr/local/poudriere/ports/lab/ports`. Describe it in
a `Makefile`, `pkg-descr`, and `pkd-plist`. Then generate a `distinfo`:

```
# cd /usr/local/poudriere/ports/lab/ports/category/example

# env PORTSDIR=/usr/local/poudriere/ports/default make makesum
```

The command above sets `PORTSDIR` to let make(1) include files from
FreeBSD-ports' `Mk/` folder.

Build the port:

```
# poudriere bulk -j 15amd64 -O lab -C category/example
```

## Port a Go application

Ref: [FreeBSD-ports/Mk/Uses/go.mk](https://github.com/freebsd/freebsd-ports/blob/0d401addc581c17e55036498065bedf462dcfb60/Mk/Uses/go.mk)

Follow standard documentation on porting Go applications to FreeBSD:

- https://docs.freebsd.org/en/books/porters-handbook/special/#using-go
- https://docs.freebsd.org/en/books/porters-handbook/uses/#uses-go

Generate `GH_TUPLE` make(1) variable for the Makefile using
https://github.com/dmgk/modules2tuple:

```
$ cd /path/to/lab

$ go mod vendor

$ modules2tuple ./vendor/modules.txt
```

Do not check-in generated `vendor/` folder. Poudriere will create it with
packages pulled from `GH_TUPLE`.

Keep in mind that Go will fail to build code with generated `vendor/` without
`vendor/modules.txt` because it runs in `-mod=vendor` mode.

Instead, store `modules.txt` in the port and inject it in the poudriere's
worktree in the Makefile:

```
post-extract:
	@${CP} ${FILESDIR}/modules.txt ${WRKSRC}/vendor
```
