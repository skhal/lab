<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**balancer** - load balancer

Synopsis
========

```console
% balancer [-n messages] [-p producers] [-w workers]
```

Description
===========

Balancer demonstrates a round robin load balancing of P producers with W workers, where each producer generates up to N messages.

It uses channels for communication to demonstrate the simplicity of writing concurrent code in Go compared to standard approaches of mutexes and other synchronization primitives in other languages, e.g. C++ or Java.

The implementation is inspired by Rob Pike's [talk ~22min](https://go.dev/blog/io2010) at Google I/O in 2010.
