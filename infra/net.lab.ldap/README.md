<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

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
