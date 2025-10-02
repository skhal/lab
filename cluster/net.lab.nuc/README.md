# NAME

**nuc.lab.net** -- Intel NUC server with FreeBSD OS


# DESCRIPTION

The hardware is Intel NUC Kit NUC6i7KYK with AMD64 architecture and the
following modules installed:

  * RAM: 2x 16 GB DDR4-2400 16GB
    HyperX HX424S14IBK2/32 kit

  * SSD: 2x 1TB NVMe
    WD WDS100T3X0C Black SN 750 Gen3 PCIe M.2 2280

## Sync files

Use `rsync(1)` to sync files with the remote host. `rsync.files-from` lists
files to pull from the server to skip history and other unnecessary files.

> [!NOTE]
> Keep `rsync.files-from` file sorted to optimize `rsync(1)`

```console
% rsync -arvz --files-from=./rsync.files-from op@nuc.lab.net:/ ./
```


# SEE ALSO

  * [Install FreeBSD](./doc/install.md)
  * [Basic Setup](./doc/setup.md)
