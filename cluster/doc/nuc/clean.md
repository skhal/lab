Name
====

**clean** - clean FreeBSD installation

Description
===========

Boot from a boot-only USB drive in a single user mode:

```console
-- destroy partitions
# gpart destroy -F /dev/nda0

-- write zeros
# dd if=/dev/zero of=/dev/nda0 bs=1m status=progress

-- (opt) write random data
# dd if=/dev/random of=/dev/nda0 bs=1m status=progress
```
