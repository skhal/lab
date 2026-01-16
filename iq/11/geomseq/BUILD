# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "find_cc",
    srcs = ["find.cc"],
    hdrs = ["find.h"],
    implementation_deps = ["@abseil-cpp//absl/strings:str_format"],
    visibility = ["//visibility:public"],
)

cc_test(
    name = "find_test",
    srcs = ["find_test.cc"],
    deps = [
        ":find_cc",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@googletest//:gtest",
        "@googletest//:gtest_main",
    ],
)
