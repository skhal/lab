# NAME

**dev.lab.net** - FreeBSD development environment


# DESCRIPTION

`dev.lab.net` is a FreeBSD development environment.

## Configuration

* tmux(1)
* vim(1)
* Go language

## Sync

```console
% rsync -arvz --files-from=./rsync.files-from op@dev.lab.net:/ ./
```


# SEE ALSO

* [Setup](./doc/setup.md)
