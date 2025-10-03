# NAME

**template-freebsd** - create a FreeBSD jail template


# DESCRIPTION

The instructions will create a ZFS dataset with patched FreeBSD user land and
ready to create a new jail. It uses FreeBSD 14.3 release for amd64 architecture.

## Bootstrap

Download FreeBSD user land into `/jail/image/`:

```console
# fetch \
    -m  \
    -o /jail/image/14.3-RELEASE-base.txz \
    https://download.freebsd.org/ftp/releases/amd64/amd64/14.3-RELEASE/base.txz
```

Open the archive into `/jail/template/` sub-folder:

```console
# tar \
    -C /jail/template/14.3-RELEASE \
    -x \
    -f /jail/image/14.3-RELEASE-base.txz \
    --unlink
```

## Patch

Update the release with latest patches. Use `PAGER=cat` to suppress
interactivity in freebsd-update(8).

```console
# env PAGER=cat freebsd-update -b /jail/template/14.3-RELEASE/ fetch install
```

Verify work:

```console
# chroot /jail/template/14.3-RELEASE freebsd-version -u
14.3-RELEASE-p2
```

## Configure

### Virtual network RC

Install [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail)
to manage jail Virtual Networks:

```console
# chroot /jail/template/14.3-RELEASE \
    fetch \
    -o /usr/local/etc/rc.jail \
    https://raw.githubusercontent.com/skhal/lab/refs/heads/main/freebsd/rc/rc.jail
```

### Root

Change shell to `tcsh(1)`:

```console
# chroot /jail/template/14.3-RELEASE \
    chsh -s /bin/tcsh root
# cp /root/.cshrc /jail/template/14.3-RELEASE/root/
```

### Routing

Set default gateway:

```console
# chroot /jail/template/14.3-RELEASE \
    sysrc -f /etc/rc.conf.d/routing defaultrouter="192.168.1.1"
# cp /etc/resolv.conf /jail/template/14.3-RELEASE/etc/resolv.conf
```

### Syslogd

Run in local mode, e.g. close sockets:

```console
# chroot /jail/template/14.3-RELEASE \
    sysrc -f /etc/rc.conf.d/syslogd syslogd_flags="-ss"
```

### Timezone

Default timezone is set to UTC. Switch to US Central Time:

```console
# chroot /jail/template/14.3-RELEASE \
    tzsetup /usr/share/zoneinfo/America/Chicago
# chroot /jail/template/14.3-RELEASE \
    tzsetup -r
```

## Snapshot

Create a ZFS snapshot for the template using the release patch number `pN` and
local changes number:

```console
# zfs snapshot zroot/jail/template/14.3-RELEASE@p2.0
```

If the template is updated in the future with local changes, say a change to
one of the RC-files, then re-snapshot the dataset with `p2.1` (assuming the
FreeBSD version is still `14.3-RELEASE-p2`).


# SEE ALSO

* https://docs.freebsd.org/en/books/handbook/jails/#creating-classic-jail
