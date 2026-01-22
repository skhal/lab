<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**lab** - R&D Lab

Description
===========

"I hear, and I forget; I see, and I remember; I do, and I understand" -- [source](https://barrypopik.com/blog/tell_me_and_i_forget_teach_me_and_i_may_remember_involve_me_and_i_will_lear)

Lab is the place to tinker with ideas, learn by doing, and don't be shy of making mistakes. This is where we grow.

Projects
--------

| Package            | Status                                                                                                                                                    | Notes                                       |
|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------|
| [`book/`](./book/) | [![Book CI](https://github.com/skhal/lab/actions/workflows/book_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/book_ci.yml)            | ideas from books                            |
| `check/`           | [![Check CI](https://github.com/skhal/lab/actions/workflows/check_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/check_ci.yml)         | [pre-commit](https://pre-commit.com) checks |
| `go/`              | [![Go CI](https://github.com/skhal/lab/actions/workflows/go_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/go_ci.yml)                  | Go libraries                                |
| `iq/`              | [![Interview Questions CI](https://github.com/skhal/lab/actions/workflows/iq_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/iq_ci.yml) | Interview Questions                         |
| [`x/`](./x/)       | [![X CI](https://github.com/skhal/lab/actions/workflows/x_ci.yml/badge.svg)](https://github.com/skhal/lab/actions/workflows/x_ci.yml)                     | experimental area                           |

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
