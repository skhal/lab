// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/23/substring/find.h"

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <ostream>
#include <string>
#include <string_view>

#include "absl/strings/str_format.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::substring {
namespace {

using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindTestParam {
  std::string_view name;
  std::string_view s;
  std::size_t n;
  std::string_view want;

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    return os << "{ s:" << std::quoted(tp.s) << " n:" << tp.n << " }";
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

TEST(FindTest, Example) {
  const std::string_view got = Find("aabcad", 2);

  EXPECT_THAT(got, Eq("aabca"));
}

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const std::string_view got = Find(tp.s, tp.n);

  EXPECT_THAT(got, Eq(tp.want));
}

const FindTestParam kFindTestParams[]{
    // n=0, no substitutions
    {.name = "empty", .s = "", .n = 0, .want = ""},
    {.name = "one letter", .s = "a", .n = 0, .want = "a"},
    {.name = "two letters distinct", .s = "ab", .n = 0, .want = "a"},
    {.name = "two letters same", .s = "aa", .n = 0, .want = "aa"},
    {.name = "three letters distinct", .s = "abc", .n = 0, .want = "a"},
    {.name = "three letters first two same", .s = "aab", .n = 0, .want = "aa"},
    {.name = "three letters first and last same",
     .s = "aba",
     .n = 0,
     .want = "a"},
    {.name = "three letters last two same", .s = "abb", .n = 0, .want = "bb"},
    {.name = "three letters same", .s = "aaa", .n = 0, .want = "aaa"},
    // n=1, one substituteion
    {.name = "empty", .s = "", .n = 1, .want = ""},
    {.name = "one letter", .s = "a", .n = 1, .want = "a"},
    {.name = "two letters distinct", .s = "ab", .n = 1, .want = "ab"},
    {.name = "two letters same", .s = "aa", .n = 1, .want = "aa"},
    {.name = "three letters distinct", .s = "abc", .n = 1, .want = "ab"},
    {.name = "three letters first two same", .s = "aab", .n = 1, .want = "aab"},
    {.name = "three letters first and last same",
     .s = "aba",
     .n = 1,
     .want = "aba"},
    {.name = "three letters last two same", .s = "abb", .n = 1, .want = "abb"},
    {.name = "three letters same", .s = "aaa", .n = 1, .want = "aaa"},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kFindTestParams),
                         [](const TestParamInfo<FindTestParam>& info) {
                           std::string name(info.param.name);
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); }, '_');
                           return absl::StrFormat("n_%d_%s", info.param.n,
                                                  name);
                         });

}  // namespace
}  // namespace iq::substring
