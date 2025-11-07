# Copyright 2025 Samvel Khalatyan. All rights reserved.

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

cc_library(
    name = "is_happy_number",
    srcs = ["is_happy_number.cc"],
    hdrs = ["is_happy_number.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "is_happy_number_test",
    srcs = ["is_happy_number_test.cc"],
    deps = [
        ":is_happy_number",
        "@abseil-cpp//absl/strings:str_format",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
