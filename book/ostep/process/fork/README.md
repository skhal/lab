<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**fork** - fork a child process

Synopsis
========

```console
./fork
```

Description
===========

A call to fork(2) creates a new child process. The child process is a fork of the parent, i.e. it continues from the fork(2) statement. Use the returned code to detect whether the code runs in the parent or child process: the child process gets return code 0.
