NAME
====

**ldap.nuc.lab.net** - OpenLDAP Directory Services server

DESCRIPTION
===========

`ldap.lab.net` hosts OpenLDAP server.

Sync
----

```console
% rsync -arvz --files-from=./rsync.files-from op@nuc.lab.net:/jail/container/ldap/ ./
```

SEE ALSO
========

-	[Setup](./doc/setup.md)
