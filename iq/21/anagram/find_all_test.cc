// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/21/anagram/find_all.h"

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <ostream>
#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::anagram {
namespace {

using ::testing::ElementsAreArray;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindAllTestParam {
  std::string name;
  std::string s;
  std::string t;
  std::vector<std::string> want;

  friend std::ostream& operator<<(std::ostream& os,
                                  const FindAllTestParam& tp) {
    return os << "s: " << std::quoted(tp.s) << " t: " << std::quoted(tp.t);
  }
};

struct FindAllTest : public TestWithParam<FindAllTestParam> {};

TEST_P(FindAllTest, Test) {
  const FindAllTestParam& tp = GetParam();

  const std::vector<std::string> got = FindAll(tp.s, tp.t);

  EXPECT_THAT(got, ElementsAreArray(tp.want));
}

const FindAllTestParam kFindAllTestParams[]{
    {.name = "empty", .s = "", .t = "", .want = {}},
    {.name = "single match", .s = "ab", .t = "ab", .want = {"ab"}},
    {.name = "single match with prefix", .s = "aab", .t = "ab", .want = {"ab"}},
    {.name = "single match with suffix", .s = "abb", .t = "ab", .want = {"ab"}},
    {.name = "two matches", .s = "aba", .t = "ab", .want = {"ab", "ba"}},
};

INSTANTIATE_TEST_SUITE_P(
    FindAllTest, FindAllTest, ValuesIn(kFindAllTestParams),
    [](const ::testing::TestParamInfo<FindAllTestParam>& info) {
      std::string name = info.param.name;
      std::replace_if(
          name.begin(), name.end(), [](char c) { return !std::isalnum(c); },
          '_');
      return name;
    });

}  // namespace
}  // namespace iq::anagram
