// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/26/cut/find.h"

#include <algorithm>
#include <cctype>
#include <ostream>
#include <string>
#include <string_view>
#include <vector>

#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "gtest/gtest.h"

namespace iq::cut {
namespace {

using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindTestParam {
  std::string_view name;
  std::vector<int> nn;
  int k;
  int want;
};

std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
  return os << absl::StrFormat("{ nn:[%s] k:%d }", absl::StrJoin(tp.nn, " "),
                               tp.k);
}

class FindTest : public TestWithParam<FindTestParam> {};

TEST(FindTest, Example) {
  const int got = Find({1, 3, 2, 4}, 2);

  EXPECT_EQ(got, 2);
}

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const int got = Find(tp.nn, tp.k);

  EXPECT_EQ(got, tp.want);
}

const FindTestParam kFindTestParams[]{
    {.name = "one item", .nn = {2}, .k = 1, .want = 1},
    {.name = "two items ascending cut one", .nn = {2, 3}, .k = 1, .want = 2},
    {.name = "two items descending cut one", .nn = {3, 2}, .k = 2, .want = 1},
    {.name = "two items ascending cut two", .nn = {2, 3}, .k = 2, .want = 1},
    {.name = "two items descending cut two", .nn = {3, 2}, .k = 2, .want = 1},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kFindTestParams),
                         [](const TestParamInfo<FindTestParam>& info) {
                           std::string name(info.param.name);
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); }, '_');
                           return name;
                         });

}  // namespace
}  // namespace iq::cut
