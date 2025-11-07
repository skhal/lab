# Copyright 2025 Samvel Khalatyan. All rights reserved.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "find_cc",
    srcs = ["find.cc"],
    hdrs = ["find.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "find_test",
    srcs = ["find_test.cc"],
    deps = [
        ":find_cc",
        "@abseil-cpp//absl/strings",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
