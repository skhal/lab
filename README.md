<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**lab** - R&D Lab

Status
======

| Package  | Status                                                                                                                                                    |
|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `book/`  | [![Book CI](https://github.com/skhal/lab/actions/workflows/book_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/book_ci.yml)            |
| `check/` | [![Check CI](https://github.com/skhal/lab/actions/workflows/check_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/check_ci.yml)         |
| `go/`    | [![Go CI](https://github.com/skhal/lab/actions/workflows/go_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/go_ci.yml)                  |
| `iq/`    | [![Interview Questions CI](https://github.com/skhal/lab/actions/workflows/iq_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/iq_ci.yml) |
| `x/`     | [![X CI](https://github.com/skhal/lab/actions/workflows/x_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/x_ci.yml)                     |

C++ development
---------------

Run `//:gazelle` target to generate Bazel `BUILD` files:

```console
% bazel run :gazelle
```

Use `refresh_commands` script to run `:gazelle` and capture `compile_commands.json` from Bazel build and update the links for FreeBSD jail installation (Ubuntu apt(1) installs LLVM under `/usr`, FreeBSD pkg(1) installs it under `/usr/local` ):

```console
% ./refresh_commands
```

[`ibazel`](https://github.com/bazelbuild/bazel-watcher) is configured to fix build errors, see [`.bazel_fix_commands.json`](./.bazel_fix_commands.json):

```console
% ibazel --run_output --run_output_interactive=false test //...
```
