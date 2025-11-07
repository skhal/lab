# Copyright 2025 Samvel Khalatyan. All rights reserved.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "last_nth_cc",
    srcs = ["last_nth.cc"],
    hdrs = ["last_nth.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "last_nth_test",
    srcs = ["last_nth_test.cc"],
    deps = [
        ":last_nth_cc",
        "@abseil-cpp//absl/strings",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
