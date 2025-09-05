# NAME

**laptop** -- Apple MacBook Air 13-inch M3 2024 with MacOS


## Sync files

Use `rsync(1)` to sync files with the remote host. `rsync.files-from` lists
files to pull from the server to skip history and other unnecessary files.

> [!NOTE]
> Keep `rsync.files-from` file sorted to optimize `rsync(1)`

```console
% rsync -arvz --files-from=./rsync.files-from / ./
```


# SEE ALSO

* [SSH Setup](./doc/ssh.md)
