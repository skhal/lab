<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**exec** - execute a different program in a given process

Synopsis
========

```console
% ./exec
[85739] start
[85739] forked child [85849]
Fri Jan 23 06:06:40 CST 2026
[85739] wait: rc 85849
```

Description
===========

A call to exec(3) loads a different program code into current process and runs it, i.e., the process does not change.
