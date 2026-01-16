# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

cc_library(
    name = "do",
    srcs = ["do.cc"],
    hdrs = ["do.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "do_test",
    srcs = ["do_test.cc"],
    deps = [
        ":do",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
