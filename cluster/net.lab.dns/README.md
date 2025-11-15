NAME
====

**dns.lab.net** - DNS server named(8)

DESCRIPTION
===========

`dns.lab.net` virtual host runs a DNS server named(8), provided by `bind9` package.

Sync
----

```console
% rsync -arvz --files-from=./rsync.files-from op@nuc.lab.net:/jail/container/dns/ ./
```
