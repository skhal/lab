# NAME

**ldap.nuc.lab.net** - LDAP server


# DESCRIPTION

`ldap.nuc.lab.net` hosts OpenLDAP server for Single Sign On on other nodes.
It is a FreeBSD jail with Virtual Network and Internet access.

## Sync files

Use `rsync(1)` to sync files with the remote host. `rsync.files-from` lists
files to pull from the server to skip history and other unnecessary files.

> [!NOTE]
> Keep `rsync.files-from` file sorted to optimize `rsync(1)`

```console
% rsync -arvz --files-from=./rsync.files-from op@dev.nuc.lab.net:/jail/container/ldap/ ./
```


# SEE ALSO

  * [Setup](./doc/setup.md)
