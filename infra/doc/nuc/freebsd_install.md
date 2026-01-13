Name
====

**install** - install FreeBSD

Create the Installation Media
=============================

-	Release: FreeBSD 15.0
-	Architecture: AMD64
-	Image Type: bootonly

See FreeBSD 15 [Announcement](https://www.freebsd.org/releases/15.0R/announce/) for the difference between image types. Bootonly image is a bootable image and requires network connection to install the OS.

```console
% curl -O https://download.freebsd.org/releases/ISO-IMAGES/15.0/FreeBSD-15.0-RELEASE-amd64-bootonly.iso.xz
% curl -O https://download.freebsd.org/releases/ISO-IMAGES/15.0/CHECKSUM.SHA512-FreeBSD-15.0-RELEASE-amd64
% shasum -c CHECKSUM.SHA512-FreeBSD-15.0-RELEASE-amd64 --ignore-missing
FreeBSD-15.0-RELEASE-amd64-bootonly.iso.xz: OK
```

Use xz(1) to decompress the file (OS X does not include xz(1), use Homebrew to install it `brew install xz`\):

```console
% xz -d ./FreeBSD-15.0-RELEASE-amd64-bootonly.iso.xz
% shasum -c CHECKSUM.SHA256-FreeBSD-15.0-RELEASE-amd64 --ignore-missing
FreeBSD-15.0-RELEASE-amd64-bootonly.iso: OK
```

Write the image to the USB drive:

```console
% disktuil list
...
/dev/disk6 (external, physical):
#:                       TYPE NAME                    SIZE       IDENTIFIER
0:     FDisk_partition_scheme                        *7.9 GB     disk6
1:                       0xEF                         34.1 MB    disk6s1
2:                    FreeBSD                         1.4 GB     disk6s2
	      (free space)                         6.5 GB     -
% sudo diskutil umountDisk /dev/disk6
Unmount of all volumes on disk6 was successful
% sudo dd status=progress if=FreeBSD-15.0-RELEASE-amd64-bootonly.iso of=/dev/disk6 bs=1m conv=sync
532676608 bytes (533 MB, 508 MiB) transferred 22.021s, 24 MB/s
508+0 records in
508+0 records out
532676608 bytes transferred in 22.059363 secs (24147416 bytes/sec)
% sudo diskutil eject /dev/disk6
Disk /dev/disk6 ejected
```

Install FreeBSD
===============

Boot from the USB drive and follow the instructions:

-	Keyboard: default layout (US)
-	Hostname: `nuc.lab.net`
-	Installation type: Packages (aka pkgbase)
-	Network: `igc0` IPv4 DHCP
-	Partition: guided ZFS
	-	Stripe with 1 disk
	-	[4g swap size](https://forums.freebsd.org/threads/swap-size-on-zfs-with-high-amount-of-ram.71059/) (increase for Kernel development)
-	Packages to install: base, lib32
-	Set root password
-	Time zone: America / USA / Central
-	Services at boot:
	-	`dumpdev` - dump kernel crashes to `/var/crash`
	-	`ntpd` - clock synchronization
	-	`ntpd_sync_on_start` - sync time on `ntpd` start
	-	`powerd` - system power control utility
	-	`sshd` - SSH server
-	Turn on security hardening:
	-	`hide_uids` hide processes as other users
	-	`hide_gids` hide processes as other groups
	-	`hide_jail` hide processes in jails
	-	`read_msgbuf` no kernel msgbuf read for unprivileged
	-	`proc_debug` no proc debug for unprivileged
	-	`random_pid` random PID for new processes
	-	`clear_tmp` clean `/tmp` on startup
	-	`disable_syslogd` no syslogd network socket
	-	`secure_console` console password prompt
-	Install firmware packages: none (skip WiFi, we'll install it later, as needed)
-	Add users: yes
	-	Name: op
	-	UID: 1000
	-	Invite to other groups: wheel
	-	Shell: tcsh

Update the OS
=============

The setup uses "Packages mechanism". It uses pkg(1) to install base system. There is no need to use freebsd-update(8) any more.

Keep an eye on the version of installed Kernel (-k), running Kernel (-r), and installed userland (-u). Reboot if the versions are out of sync after an upgrade of the packages:

```console
% freebsd-version -kru
15.0-RELEASE
15.0-RELEASE
15.0-RELEASE
```

The list of intalled packages is minimal. It includes the kernel, minimal packages set to run multi-user environment, base userland, and 32-bit compatibility libraries:

```console
% pkg prime-list
FreeBSD-kernel-generic
FreeBSD-set-base
FreeBSD-set-lib32
FreeBSD-set-minimal
pkg
```

See Also
========

-	[FreeBSD Handbook - Install](https://docs.freebsd.org/en/books/handbook/bsdinstall/)
