# Bootstrap

The following instructions demonstrate how to install FreeBSD 14.3 release.


## Step 1: Create a bootable disk.

Ref: https://docs.freebsd.org/en/books/handbook/bsdinstall/

* Download `amd64` ISO image from FreeBSD
  [download server](https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/).

  ```
  % curl \
      -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/FreeBSD-14.3-RELEASE-amd64-bootonly.iso
  ```

  > **Note**: If using Mac to write the image, prefer `.iso` file unless
  > `xz-utils` are installed. Even though `tar -xjf` does recognize XZ format
  > it extracts the contents of the archive rather than keeping the ISO image.

* Verify the checksum file from the same server.

  ```
  % curl \
      -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64

  % sha512sum \
      -c CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64 \
      FreeBSD-14.3-RELEASE-amd64-bootonly.iso
  FreeBSD-14.3-RELEASE-amd64-bootonly.iso: OK
  ...
  ```

* To write ISO image, one needs to identify external memory stick. Connect the
  external card and verify it is available in `/dev` using `diskutil`. Then
  unmount it and write the image:

  ```
  % disktuil list
  ...
  /dev/disk6 (external, physical):
     #:                       TYPE NAME                    SIZE       IDENTIFIER
     0:     FDisk_partition_scheme                        *7.9 GB     disk6
     1:                       0xEF                         34.1 MB    disk6s1
     2:                    FreeBSD                         1.4 GB     disk6s2
                      (free space)                         6.5 GB     -

  % diskutil umountDisk /dev/disk6
  Unmount of all volumes on disk6 was successful

  % sudo dd \
      status=progress \
      if=FreeBSD-14.3-RELEASE-amd64-bootonly.iso \
      of=/dev/disk6 \
      bs=1M \
      conv=sync
  407896064 bytes (408 MB, 389 MiB) transferred 17.011s, 24 MB/s
  412+0 records in
  412+0 records out
  432013312 bytes transferred in 18.013281 secs (23983044 bytes/sec)

  % diskutil eject /dev/disk6
  Disk /dev/disk6 ejected
  ```
