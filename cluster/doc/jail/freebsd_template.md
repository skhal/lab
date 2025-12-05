Name
====

**freebsd-template** - create a FreeBSD jail template

Bootstrap
=========

Download FreeBSD userland into `/jail/image/`:

```console
# fetch -m -o /jail/image/15.0-RELEASE-base.txz https://download.freebsd.org/ftp/releases/amd64/15.0-RELEASE/base.txz
```

Extract the archive into a new ZFS dataset:

```console
# zfs create zroot/jail/template/15.0-RELEASE
# tar -C /jail/template/15.0-RELEASE -x -f /jail/image/15.0-RELEASE-base.txz --unlink
```

Patch
=====

Update the release with latest patches. Use `PAGER=cat` to suppress interactivity in freebsd-update(8).

```console
# env PAGER=cat freebsd-update -b /jail/template/15.0-RELEASE/ fetch install
```

Verify work:

```console
# chroot /jail/template/15.0-RELEASE/ freebsd-version -u
15.0-RELEASE
```

Configure
---------

Change root shell to `tcsh(1)`:

```console
# chroot /jail/template/15.0-RELEASE/ chsh -s /bin/tcsh root
# cp /root/.cshrc /jail/template/15.0-RELEASE/root/
```

Routing
-------

Set default gateway:

```console
# chroot /jail/template/15.0-RELEASE/ sysrc -f /etc/rc.conf.d/routing defaultrouter="192.168.1.1"
cp /etc/resolv.conf /jail/template/15.0-RELEASE/etc/
```

Syslogd
-------

Run in local mode to close listening socket:

```console
# chroot /jail/template/15.0-RELEASE/ sysrc -f /etc/rc.conf.d/syslogd syslogd_flags="-ss"
```

### Timezone

Default timezone is set to UTC. Switch to US Central Time:

```console
# chroot /jail/template/15.0-RELEASE/ tzsetup America/Chicago
# chroot /jail/template/15.0-RELEASE/ tzsetup -r
```

Snapshot
--------

Create a ZFS snapshot for the template using the release patch number `pN` and local changes number:

```console
# zfs snapshot zroot/jail/template/15.0-RELEASE@p0.0
```

SEE ALSO
========

-	https://docs.freebsd.org/en/books/handbook/jails/#creating-classic-jail
