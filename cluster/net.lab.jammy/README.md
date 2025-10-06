# NAME

**jammy.lab.net** - Ubuntu 22.04 development environment


# DESCRIPTION

`jammy.lab.net` is a Ubuntu 22.04 development environment.

## Configuration

* tmux(1)
* vim(1)
* Bazel 8.4.2

## Sync

```console
% rsync -arvz --files-from=./rsync.files-from op@nuc.lab.net:/jail/container/jammy/ ./
```


# SEE ALSO

* [Setup](./doc/setup.md)
