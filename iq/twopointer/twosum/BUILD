# Copyright 2025 Samvel Khalatyan. All rights reserved.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "solution_cc",
    srcs = ["solution.cc"],
    hdrs = ["solution.h"],
)

cc_test(
    name = "solution_test",
    srcs = ["solution_test.cc"],
    deps = [
        ":solution_cc",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
