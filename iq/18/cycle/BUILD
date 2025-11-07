# Copyright 2025 Samvel Khalatyan. All rights reserved.

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
        "@googletest//:gtest_main",
    ],
)
