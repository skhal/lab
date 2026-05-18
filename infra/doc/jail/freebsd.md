<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**freebsd** - FreeBSD environment in a jail

# DESCRIPTION

The instructions create a FreeBSD jail template with network setup, that is
ready for cloning to spawn new jails.

## Template

ZFS dataset:

```console
# setenv OS_VERSION=$(freebsd-version | cut -d'-' -f 1)
# echo $OS_VERSION
15.0

# zfs create -v zroot/usr/jail/${OS_VERSION}
```

Install minimal environment:

```console
-- pkg(8) keys to access repos
# mkdir -vp /usr/jail/${OS_VERSION}/usr/share/keys
# cp -vrn /usr/share/keys/pkg* /usr/jail/${OS_VERSION}/usr/share/keys/

-- hierarchy for pkg(8) database
# mkdir -vp /usr/jail/${OS_VERSION}/var/db/pkg/repos

-- metadata cache
# mount -vt tmpfs tmpfs /usr/jail/${OS_VERSION}/var/db/pkg/repos

# pkg -r /usr/jail/${OS_VERSION} install -y FreeBSD-set-minimal-jail pkg

# chroot /usr/jail/${OS_VERSION} freebsd-version -u
15.0-RELEASE
```

Root shell and password:

```
# pkg -r /usr/jail/${OS_VERSION}/ install -y FreeBSD-csh neovim
# chroot /usr/jail/${OS_VERSION}/ chsh -s /bin/tcsh root
# cp /root/.cshrc /usr/jail/${OS_VERSION}/root/

# chroot /usr/jail/${OS_VERSION}/ passwd
```

Network and services (change default router to jailed DNS when ready):

```
# pkg -r /usr/jail/${OS_VERSION}/ install FreeBSD-bsdconfig

# chroot /usr/jail/${OS_VERSION}/ sysrc -f /etc/rc.conf.d/routing defaultrouter="192.168.1.1"
# cp /etc/resolv.conf /usr/jail/${OS_VERSION}/etc/

-- syslogd in secure mode, no network socket, no compression (ZFS compresses)
# chroot /usr/jail/${OS_VERSION}/ sysrc -f /etc/rc.conf.d/syslogd syslogd_flags="-c -ss"

-- UTC to US Central Time
# chroot /usr/jail/${OS_VERSION}/ tzsetup America/Chicago
# chroot /usr/jail/${OS_VERSION}/ tzsetup -r
```

Enable FreeBSD repos with final upgrade:

```console
# mkdir -p /usr/jail/${OS_VERSION}/usr/local/etc/pkg/repos

# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/jail/data/template/usr/local/etc/pkg/repos/FreeBSD.conf
# mv -nv /tmp/FreeBSD.conf /usr/jail/${OS_VERSION}/usr/local/etc/pkg/repos/FreeBSD.conf

# pkg update
# pkg upgrade -y

# zfs snapshot zroot/usr/jail/15.0@`date +%y%m%d`
```
