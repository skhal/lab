# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

cc_library(
    name = "find_all",
    srcs = ["find_all.cc"],
    hdrs = ["find_all.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "find_all_test",
    srcs = ["find_all_test.cc"],
    deps = [
        ":find_all",
        "@googletest//:gtest_main",
    ],
)
