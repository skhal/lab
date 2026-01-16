// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/6/lexiseq/solution.h"

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <ostream>
#include <vector>

#include "gtest/gtest.h"

namespace iq::lexiseq {
namespace {

using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

struct NextTestParam {
  std::string name;
  std::string s;
  std::string want;

  friend std::ostream& operator<<(std::ostream& os, const NextTestParam& tp) {
    return os << std::quoted(tp.name);
  }
};

class NextTest : public TestWithParam<NextTestParam> {};

TEST_P(NextTest, Test) {
  const NextTestParam& tp = GetParam();

  const std::string got = Next(tp.s);

  EXPECT_EQ(got, tp.want);
}

const std::vector<NextTestParam> kTestParams{
    {.name = "empty", .s = "", .want = ""},
    {.name = "one letter", .s = "a", .want = "a"},
    {.name = "two letters same", .s = "aa", .want = "aa"},
    {.name = "two letters ascending", .s = "ab", .want = "ba"},
    {.name = "two letters descending", .s = "ba", .want = "ab"},
    {.name = "three letters same", .s = "aaa", .want = "aaa"},
    {.name = "three letters first two same", .s = "aab", .want = "aba"},
    {.name = "three letters first and last same", .s = "aba", .want = "baa"},
    {.name = "three letters second and last same", .s = "baa", .want = "aab"},
};

INSTANTIATE_TEST_SUITE_P(NextTest, NextTest, ValuesIn(kTestParams),
                         [](const TestParamInfo<NextTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::lexiseq
