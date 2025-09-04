# NAME

**install** - install FreeBSD


# DESCRIPTION

Ref: https://docs.freebsd.org/en/books/handbook/bsdinstall/

## Create the Installation Media

*Note*: The instructions use MacOS to create the installation media. The OS does not
include include `xz(1)` to uncompress images from `.xg` files.

Download AMD64 image and verify the checksum.

```console
% curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/FreeBSD-14.3-RELEASE-amd64-mini-memstick.img
% curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64
% shasum \
    -c CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64 \
    --ignore-missing
FreeBSD-14.3-RELEASE-amd64-mini-memstick.img: OK
```

Write the image to the memstick:

```console
% disktuil list
...
/dev/disk6 (external, physical):
#:                       TYPE NAME                    SIZE       IDENTIFIER
0:     FDisk_partition_scheme                        *7.9 GB     disk6
1:                       0xEF                         34.1 MB    disk6s1
2:                    FreeBSD                         1.4 GB     disk6s2
	      (free space)                         6.5 GB     -
# diskutil umountDisk /dev/disk6
Unmount of all volumes on disk6 was successful
# dd \
    status=progress \
    if=FreeBSD-14.3-RELEASE-amd64-mini-memstick.img \
    of=/dev/disk6 \
    bs=1m \
    conv=sync
532676608 bytes (533 MB, 508 MiB) transferred 22.021s, 24 MB/s
508+0 records in
508+0 records out
532676608 bytes transferred in 22.059363 secs (24147416 bytes/sec)
# diskutil eject /dev/disk6
Disk /dev/disk6 ejected
```

## Install the OS

Boot from the memstick and follow the instructions:

  * Keyboard: default layout (US)
  * Hostname: `nuc.lab.net`
  * Network: `em0` IPv4 DHCP
  * Filesystem: ZFS with mirrored pool
  * Enable services:
    - `dumpdev` to dump kernel crashes to `/var/crash`
    - `ntpd` for clock synchronization
    - `ntpd_sync_on_start` to sync time on `ntpd` start
    - `sshd` for SSH server
  * Turn on security hardening:
    - `hide_uids` hide processes as other users
    - `hide_gids` hide processes as other groups
    - `hide_jail` hide processes in jails
    - `read_msgbuf` no kernel msgbuf read for unprivileged
    - `proc_debug` no proc debug for unprivileged
    - `random_pid` random PID for new processes
    - `clear_tmp` clean `/tmp` on startup
    - `disable_syslogd` no syslogd network socket
    - `secure_console` console password prompt

## Update the OS

Ref: https://docs.freebsd.org/en/books/handbook/cutting-edge/

Check the running version of installed kernel `-k`, running kernel `-r`,
and running userland `-u`:

```console
% freebsd-version -kru
14.3-RELEASE
14.3-RELEASE
14.3-RELEASE
```

Fetch and install updates. The system will auto-reboot if there is
a kernel update, otherwise it restarts the updated services only.
It is still a good idea to restart the node.

```console
# freebsd-update fetch
# freebsd-update install
# reboot
```

Verify the running version is up to date, all three should match:

```console
% freebsd-version -kru
14.3-RELEASE-p2
14.3-RELEASE-p2
14.3-RELEASE-p2
```
