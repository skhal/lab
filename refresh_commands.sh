#!/bin/sh
#
# Copyright 2025 Samvel Khalatyan. All rights reserved.

bazel run :refresh_compile_commands

sed \
  -i'' \
  -e 's,/usr/bin/clang-21,/usr/local/bin/clang21,' \
  -e 's,/usr/lib/llvm-21,/usr/local/llvm21,' \
  compile_commands.json
