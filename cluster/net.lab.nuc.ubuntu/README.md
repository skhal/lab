# NAME

**ubuntu.nuc.lab.net** - Ubuntu jail on `nuc.lab.net` server


# DESCRIPTION

`ubuntu.nuc.lab.net` is a vnet jail with Internet access, running Ubuntu on
`nuc.lab.net` using Linux API and bootstrapped Ubuntu libraries with
`debootstrap(8)`.

## Sync files

Use `rsync(1)` to sync files with the remote host. `rsync.files-from` lists
files to pull from the server to skip history and other unnecessary files.

> [!NOTE]
> Keep `rsync.files-from` file sorted to optimize `rsync(1)`

```console
% rsync -arvz --files-from=./rsync.files-from op@nuc.lab.net:/jail/container/ubuntu/ ./
```


# SEE ALSO

  * [Basic Setup](./doc/setup.md)
