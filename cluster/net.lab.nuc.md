# Bootstrap

The following instructions demonstrate how to install FreeBSD 14.3 release.


## Step 1: Create a bootable memstick.

Ref: https://docs.freebsd.org/en/books/handbook/bsdinstall/

The instructions download `.img` for memstick. `.iso` file is for optical media.
Even though there are `xz(1)` compressed images, MacOS does not include the tool
to uncompress these. As such, the steps below download `.img`.

* Download `amd64` ISO image, checksum list, and verify download
  by matching checksum:

  ```console
  % curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/FreeBSD-14.3-RELEASE-amd64-mini-memstick.img

  % curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64

  % shasum -c CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64 --ignore-missing
  FreeBSD-14.3-RELEASE-amd64-mini-memstick.img: OK
  ```

* Insert memstick, identify the device and prepare it to write the image using
  `diskutil(1)`.

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
  ```

* Write the image:

  ```console
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


## Step 3: Update the system

Ref: https://docs.freebsd.org/en/books/handbook/cutting-edge/

* Check running version of installed kernel `-k`, running kernel `-r`, and
  running userland `-u`:

  ```console
  % freebsd-version -kru
  14.3-RELEASE
  14.3-RELEASE
  14.3-RELEASE
  ```

  alternatively:

  ```
  % uname -a
  FreeBSD nuc.lab.net 14.3-RELEASE FreeBSD 14.3-RELEASE GENERIC amd64
  ```

* Fetch and install updates--the system will auto-reboot if there is
  a kernel update, otherwise it restarts updated services only. It is
  still a good idea to restart the node.

  ```
  # freebsd-update fetch
  # freebsd-update install
  # reboot
  ```

* Verify running version is up to date, e.g. the installed and running
  running kernels should match, the userland should match these too:

  ```console
  % freebsd-version -kru
  14.3-RELEASE-p2
  14.3-RELEASE-p2
  14.3-RELEASE-p2
  ```


## Step 4: Setup

Create an operator user:

* [one time] create an template configuration `/etc/adduser.conf` for
  `adduser(1)` to use `tcsh(1)` by default:

  ```console
  # adduser -C
  ...
  Shell (sh csh tcsh zsh nologin) [sh]: tcsh
  ...
  ```

* Create an operator group `op:1000` and user `op:1000` to manage the host:

  ```console
  # pw groupadd -n op -g 1000
  # pw useradd \
      -c 'System Operator' \
      -d /home/op \
      -g op \
      -G wheel \
      -m \
      -n op \
      -s /bin/tcsh \
      -u 1000 \
      -w no
  # passwd -l op
  ```
