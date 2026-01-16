// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/3/palindrome/solution.h"

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <ostream>
#include <string>

#include "gtest/gtest.h"

namespace iq::palindrome {
namespace {

using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

struct IsTestParam {
  std::string name;
  std::string s;
  bool want;

  friend std::ostream& operator<<(std::ostream& os, const IsTestParam& tp) {
    return os << std::quoted(tp.name) << " s: " << std::quoted(tp.s);
  }
};

class IsTest : public TestWithParam<IsTestParam> {};

TEST_P(IsTest, pass) {
  const IsTestParam& tp = GetParam();

  const bool got = Is(tp.s);

  ASSERT_EQ(got, tp.want);
}

const IsTestParam kIsTestParams[]{
    {.name = "empty", .s = "", .want = true},
    {.name = "one letter", .s = "a", .want = true},
    {.name = "one digit", .s = "1", .want = true},
    {.name = "one non-alphanumeric letter", .s = ".", .want = true},
    {.name = "two letters", .s = "aa", .want = true},
    {.name = "two digits", .s = "11", .want = true},
    {.name = "two distinct digits", .s = "12", .want = false},
    {.name = "letter and digit", .s = "a1", .want = false},
    {.name = "letter and non-alphanumeric", .s = "a.", .want = true},
    {.name = "digit and non-alphanumeric", .s = "1.", .want = true},
    {.name = "non-alphanumeric and letter", .s = ".a", .want = true},
    {.name = "non-alphanumeric and digit", .s = ".1", .want = true},
};

INSTANTIATE_TEST_SUITE_P(IsTest, IsTest, ValuesIn(kIsTestParams),
                         [](const TestParamInfo<IsTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::palindrome
