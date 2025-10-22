# Copyright 2025 Samvel Khalatyan. All rights reserved.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "clean_cc",
    srcs = ["clean.cc"],
    hdrs = ["clean.h"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "clean_test",
    srcs = ["clean_test.cc"],
    deps = [
        ":clean_cc",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
