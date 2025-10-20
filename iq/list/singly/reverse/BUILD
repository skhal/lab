# Copyright 2025 Samvel Khalatyan. All rights reserved.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "reverse_cc",
    srcs = ["reverse.cc"],
    hdrs = ["reverse.h"],
)

cc_test(
    name = "reverse_test",
    srcs = ["reverse_test.cc"],
    deps = [
        ":reverse_cc",
        "@abseil-cpp//absl/strings",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
