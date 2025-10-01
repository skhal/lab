cc_library(
  name = "solution_cc",
  hdrs = ["solution.h"],
  srcs = ["solution.cc"],
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
