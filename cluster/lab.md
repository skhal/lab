# Bootstrap

The following instructions demonstrate how to install FreeBSD 14.3 release.


## Step 1: Create a bootable memstick.

Ref: https://docs.freebsd.org/en/books/handbook/bsdinstall/

> [!NOTE]
> FreeBSD download server includes `.xz` compressed images to save space.
> As of Aug 2025, MacOS does not include `xz(1)`. `tar(1)` supports XZ
> compressed archives--it extracts contents of the archive instead of
> keeping the image `.iso` or `.img`, which is to be written to the external
> memstick with `dd(1)`. Therefore the instructions below download `.img`
> instead of `.img.xz`.

* Download `amd64` ISO image from FreeBSD
  [download server](https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/).

  ```console
  % curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/FreeBSD-14.3-RELEASE-amd64-mini-memstick.img
  ```

* Download checksum list:

  ```console
  % curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64
  ```

* Verify checksum:

  ```console
  % shasum -c CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64 --ignore-missing
  FreeBSD-14.3-RELEASE-amd64-mini-memstick.img: OK
  ```

* Use `diskutil(1)` to identify the memstick to write the image to, assuming
  it is connected to laptop.

  ```console
  % disktuil list
  ...
  /dev/disk6 (external, physical):
     #:                       TYPE NAME                    SIZE       IDENTIFIER
     0:     FDisk_partition_scheme                        *7.9 GB     disk6
     1:                       0xEF                         34.1 MB    disk6s1
     2:                    FreeBSD                         1.4 GB     disk6s2
                      (free space)                         6.5 GB     -
  ```

  Unmount memstick:

  ```console
  % diskutil umountDisk /dev/disk6
  Unmount of all volumes on disk6 was successful
  ```

  Write the image:

  ```console
  % sudo dd \
      status=progress \
      if=FreeBSD-14.3-RELEASE-amd64-mini-memstick.img \
      of=/dev/disk6 \
      bs=1m \
      conv=sync
  532676608 bytes (533 MB, 508 MiB) transferred 22.021s, 24 MB/s
  508+0 records in
  508+0 records out
  532676608 bytes transferred in 22.059363 secs (24147416 bytes/sec)
  ```

  Eject the memstick:

  ```console
  % diskutil eject /dev/disk6
  Disk /dev/disk6 ejected
  ```

## Step 2: Install

Boot from the memstick and follow the instructions:

* Keyboard: default layout (US)
* Hostname: `nuc.lab.net`
* Network: `em0` IPv4 DHCP
* Filesystem: ZFS with mirrored pool
* Services:
  - `dumpdev` to dump kernel crashes to `/var/crash`
  - `ntpd` for clock synchronization
  - `ntpd_sync_on_start` to sync time on `ntpd` start
  - `sshd` for SSH server
* Security hardening:
  - `hide_uids` hide processes as other users
  - `hide_gids` hide processes as other groups
  - `hide_jail` hide processes in jails
  - `read_msgbuf` no kernel msgbuf read for unprivileged
  - `proc_debug` no proc debug for unprivileged
  - `random_pid` random PID for new processes
  - `clear_tmp` clean `/tmp` on startup
  - `disable_syslogd` no syslogd network socket
  - `secure_console` console password prompt
* Install FreeBSD handbook

> [!NOTE]
> There are no users setup at this point except `root`.
