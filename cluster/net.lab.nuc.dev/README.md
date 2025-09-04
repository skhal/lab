# NAME

**dev.nuc.lab.net** - development jail on `nuc.lab.net` server


# DESCRIPTION

`dev.nuc.lab.net` is a vnet jail with Internet access, running on `nuc.lab.net`
FreeBSD server.

## Sync files

Use `rsync(1)` to sync files with the remote host. `rsync.files-from` lists
files to pull from the server to skip history and other unnecessary files.

> [!NOTE]
> Keep `rsync.files-from` file sorted to optimize `rsync(1)`

```console
% rsync -arvz --files-from=./rsync.files-from op@dev.nuc.lab.net:/ ./
```


# SEE ALSO

  * [Basic Setup](./doc/setup.md)
