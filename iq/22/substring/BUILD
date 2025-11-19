# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

cc_library(
    name = "find_fast",
    srcs = ["find_fast.cc"],
    hdrs = ["find_fast.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "find_fast_test",
    srcs = ["find_fast_test.cc"],
    deps = [
        ":find_fast",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
