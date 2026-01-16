<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**dev.lab.net** - FreeBSD development environment

DESCRIPTION
===========

`dev.lab.net` is a FreeBSD development environment.

Configuration
-------------

-	tmux(1)
-	vim(1)
-	Go language

Sync
----

```console
% rsync -arvz --files-from=./rsync.files-from op@dev.lab.net:/ ./
```

SEE ALSO
========

-	[Setup](./doc/setup.md)
