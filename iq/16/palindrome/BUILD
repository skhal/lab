# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

cc_library(
    name = "is",
    srcs = ["is.cc"],
    hdrs = ["is.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "is_test",
    srcs = ["is_test.cc"],
    deps = [
        ":is",
        "@abseil-cpp//absl/strings",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
