# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")
load(
    "@hedron_compile_commands//:refresh_compile_commands.bzl",
    "refresh_compile_commands",
)

gazelle_binary(
    name = "gazelle_cc",
    languages = [
        "@gazelle_cc//language/cc",
    ],
)

gazelle(
    name = "gazelle",
    gazelle = ":gazelle_cc",
)

# keep-sorted start
# gazelle:cc_group unit
# gazelle:exclude cluster
# gazelle:exclude nvim
# gazelle:exclude freebsd
# gazelle:exclude home
# gazelle:exclude toolchain
# gazelle:exclude vim
# gazelle:resolve cc gtest/gtest.h @googletest//:gtest_main
# keep-sorted end

# https://github.com/hedronvision/bazel-compile-commands-extractor?tab=readme-ov-file
refresh_compile_commands(
    name = "refresh_compile_commands",
    targets = [
        "//iq/...",
        "//x/abseil/...",
    ],
)
