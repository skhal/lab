Name
====

**freebsd-template** - create a FreeBSD jail with freebsd-base(7)

Bootstrap
=========

Prepare a new ZFS dataset:

```console
# zfs create zroot/jail/template/$(freebsd-version)
```

Copy pkg(8) keys to access repositories:

```console
# mkdir -vp /jail/template/$(freebsd-version)/usr/share/keys
# cp -vrn /usr/share/keys/pkg* /jail/template/$(freebsd-version)/usr/share/keys/
```

Prepare file hierarchy for pkg(8) database:

```console
# mkdir -vp /jail/template/$(freebsd-version)/var/db/pkg/repos
```

Mount tmpfs(4) for metadata cache:

```console
# mount -vt tmpfs tmpfs /jail/template/$(freebsd-version)/var/db/pkg/repos
```

Install minimal set of packages designed for jails:

```console
# pkg -r /jail/template/15.0-minimal install FreeBSD-set-minimal-jail
```

Verify work:

```console
# chroot /jail/template/$(freebsd-version)/ freebsd-version -u
15.0-RELEASE
```

IMPORTANT: 32bit c-libs
-----------------------

Go static linking needs to access `/usr/bin/cc`, provided by `FreeBSD-clang-15.0` package, which depends on `FreeBSD-clibs-lib32-15.0` for a bunch of libraries including `libc`:

```console
% pkg which /usr/bin/cc
/usr/bin/cc was installed by package FreeBSD-clang-15.0
% pkg info -d FreeBSD-clang-15.0
FreeBSD-clang-15.0:
        FreeBSD-lld-15.0
        FreeBSD-libcompiler_rt-dev-15.0
        FreeBSD-clibs-15.0 (libc++.so.1)
        FreeBSD-clibs-15.0 (libc.so.7)
        FreeBSD-clibs-lib32-15.0 (libc.so.7:32)
        FreeBSD-clibs-15.0 (libcxxrt.so.1)
        FreeBSD-libexecinfo-15.0 (libexecinfo.so.1)
        FreeBSD-clibs-15.0 (libgcc_s.so.1)
        FreeBSD-clibs-lib32-15.0 (libgcc_s.so.1:32)
        FreeBSD-clibs-15.0 (libm.so.5)
        FreeBSD-ncurses-lib-15.0 (libncursesw.so.9)
        FreeBSD-runtime-15.0 (libprivatezstd.so.5)
        FreeBSD-clibs-15.0 (libthr.so.3)
        FreeBSD-ncurses-lib-15.0 (libtinfow.so.9)
        FreeBSD-zlib-15.0 (libz.so.6)
```

The package `FreeBSD-clibs-lib32-15.0` installs `/libexec/ld-elf32.so.1`:

```console
% pkg info -l FreeBSD-clibs-lib32-15.0
FreeBSD-clibs-lib32-15.0:
        /libexec/ld-elf32.so.1
        /usr/lib32/libc++.so.1
        /usr/lib32/libc.so.7
        /usr/lib32/libcxxrt.so.1
        /usr/lib32/libdl.so.1
        /usr/lib32/libgcc_s.so.1
        /usr/lib32/libm.so.5
        /usr/lib32/librt.so.1
        /usr/lib32/libssp.so.0
        /usr/lib32/libsys.so.7
        /usr/lib32/libthr.so.3
        /usr/lib32/libxnet.so
        /usr/libexec/ld-elf32.so.1
```

The library has chflags(2) set:

```console
% ls -lo /libexec/ld-elf32.so.1
-r-xr-xr-x  1 root wheel schg 104792 Nov 27 18:00 /libexec/ld-elf32.so.1
```

Jails can't modify flags even by privileged users unless `allow.chflags` is enabled in the jail configuration, jail(8).

A workaround this problem is to install the above mentioned package from the host system in order to keep `allow.chflags` disabled for security reasons:

```console
# pkg -r /jail/template/$(freebsd-version)/ install FreeBSD-clibs-lib32-15.0
```

Configure
=========

Root
----

Install tcsh(1):

```console
# pkg -r /jail/template/$(freebsd-version)/ install FreeBSD-csh
```

Change root shell to tcsh(1):

```console
# chroot /jail/template/$(freebsd-version)/ chsh -s /bin/tcsh root
# cp /root/.cshrc /jail/template/$(freebsd-version)/root/
```

Routing
-------

Install sysrc(8):

```console
# pkg -r /jail/template/$(freebsd-version)/ install FreeBSD-bsdconfig
```

Set default gateway:

```console
# chroot /jail/template/$(freebsd-version)/ sysrc -f /etc/rc.conf.d/routing defaultrouter="192.168.1.1"
# cp /etc/resolv.conf /jail/template/$(freebsd-version)/etc/
```

Syslogd
-------

Run in local mode to close listening socket:

```console
# chroot /jail/template/$(freebsd-version)/ sysrc -f /etc/rc.conf.d/syslogd syslogd_flags="-ss"
```

Timezone
--------

Default timezone is set to UTC. Switch to US Central Time:

```console
# chroot /jail/template/$(freebsd-version)/ tzsetup America/Chicago
# chroot /jail/template/$(freebsd-version)/ tzsetup -r
```

Snapshot
--------

Create a ZFS snapshot for the template using the release patch number `pN` and local changes number:

```console
# zfs snapshot zroot/jail/template/$(freebsd-version)@p0.0
```

References
==========

-	https://docs.freebsd.org/en/books/handbook/jails/#creating-classic-jail
