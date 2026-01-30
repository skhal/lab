// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/22/substring/find_fast.h"

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <ostream>
#include <string_view>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::substring {
namespace {

using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindFastTestParam {
  std::string_view name;
  std::string_view s;
  std::string_view want;

  friend std::ostream& operator<<(std::ostream& os,
                                  const FindFastTestParam& tp) {
    return os << "s: " << std::quoted(tp.s);
  }
};

class FindFastTest : public TestWithParam<FindFastTestParam> {};

TEST_P(FindFastTest, Test) {
  const FindFastTestParam& tp = GetParam();

  const std::string_view got = FindFast(tp.s);

  EXPECT_THAT(got, Eq(tp.want));
}

const FindFastTestParam kFindFastTestParams[]{
    {.name = "empty", .s = "", .want = ""},
    {.name = "len one", .s = "a", .want = "a"},
    {.name = "len two", .s = "ab", .want = "ab"},
    {.name = "len two same char", .s = "aa", .want = "a"},
    {.name = "len three", .s = "abc", .want = "abc"},
    {.name = "len three same char", .s = "aaa", .want = "a"},
    {.name = "len three first two chars", .s = "aba", .want = "ab"},
};

INSTANTIATE_TEST_SUITE_P(FindFastTest, FindFastTest,
                         ValuesIn(kFindFastTestParams),
                         [](const TestParamInfo<FindFastTestParam>& info) {
                           std::string name = std::string(info.param.name);
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); }, '_');
                           return name;
                         });

}  // namespace
}  // namespace iq::substring
