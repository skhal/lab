// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/29/rotated/find.h"

#include <algorithm>
#include <cctype>
#include <ostream>
#include <string>
#include <string_view>
#include <vector>

#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::rotated {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::IsOkAndHolds;
using ::testing::Not;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindTestParam {
  std::string_view name;
  std::vector<int> nn;
  int k;
  int want;
};

std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
  return os << absl::StrFormat("{nn:[%s] k:%d}", absl::StrJoin(tp.nn, " "),
                               tp.k);
}

auto NameTest = [](const ::testing::TestParamInfo<FindTestParam>& info) {
  std::string name(info.param.name);
  std::replace_if(
      name.begin(), name.end(), [](char c) { return !std::isalnum(c); }, '_');
  return name;
};

class FindTest : public TestWithParam<FindTestParam> {};

using PositiveFindTest = FindTest;

TEST_P(PositiveFindTest, Test) {
  const FindTestParam& tp = GetParam();

  const absl::StatusOr<int> got = Find(tp.nn, tp.k);

  EXPECT_THAT(got, IsOkAndHolds(tp.want));
}

const FindTestParam kPositiveFindTestParams[]{
    // size 1
    {.name = "size one hit", .nn = {1}, .k = 1, .want = 0},
    // size 2 - no rotation
    {.name = "size two no rotation hit first", .nn = {1, 3}, .k = 1, .want = 0},
    {.name = "size two no rotation hit second",
     .nn = {1, 3},
     .k = 3,
     .want = 1},
    // size 2 - rotate one
    {.name = "size two rotate one hit first", .nn = {3, 1}, .k = 1, .want = 1},
    {.name = "size two rotate one hit second", .nn = {3, 1}, .k = 3, .want = 0},
    // size 3 - no rotation
    {.name = "size three no rotation hit first",
     .nn = {1, 3, 5},
     .k = 1,
     .want = 0},
    {.name = "size three no rotation hit second",
     .nn = {1, 3, 5},
     .k = 3,
     .want = 1},
    {.name = "size three no rotation hit third",
     .nn = {1, 3, 5},
     .k = 5,
     .want = 2},
    // size 3 - rotate one
    {.name = "size three rotate one hit first",
     .nn = {3, 5, 1},
     .k = 3,
     .want = 0},
    {.name = "size three rotate one hit second",
     .nn = {3, 5, 1},
     .k = 5,
     .want = 1},
    {.name = "size three rotate one hit third",
     .nn = {3, 5, 1},
     .k = 1,
     .want = 2},
    // size 3 - rotate two
    {.name = "size three rotate two hit first",
     .nn = {5, 1, 3},
     .k = 5,
     .want = 0},
    {.name = "size three rotate two hit second",
     .nn = {5, 1, 3},
     .k = 1,
     .want = 1},
    {.name = "size three rotate two hit third",
     .nn = {5, 1, 3},
     .k = 3,
     .want = 2},
};

INSTANTIATE_TEST_SUITE_P(FindTest, PositiveFindTest,
                         ValuesIn(kPositiveFindTestParams), NameTest);

using NegativeFindTest = FindTest;

TEST_P(NegativeFindTest, Test) {
  const FindTestParam& tp = GetParam();

  const absl::StatusOr<int> got = Find(tp.nn, tp.k);

  EXPECT_THAT(got, Not(IsOk()));
}

inline constexpr int kIllegalIndex = -1;

const FindTestParam kNegativeFindTestParams[]{
    {.name = "empty", .nn = {}, .k = 1, .want = kIllegalIndex},
    // size 1
    {.name = "size one miss below", .nn = {1}, .k = 0, .want = kIllegalIndex},
    {.name = "size one miss above", .nn = {1}, .k = 2, .want = kIllegalIndex},
    // size 2 - no rotation
    {.name = "size two no rotation miss below",
     .nn = {1, 3},
     .k = 0,
     .want = kIllegalIndex},
    {.name = "size two no rotation miss between",
     .nn = {1, 3},
     .k = 2,
     .want = kIllegalIndex},
    {.name = "size two no rotation miss above",
     .nn = {1, 3},
     .k = 4,
     .want = kIllegalIndex},
    // size 2 - rotate one
    {.name = "size two rotate one miss below",
     .nn = {3, 1},
     .k = 0,
     .want = kIllegalIndex},
    {.name = "size two rotate one miss between",
     .nn = {3, 1},
     .k = 2,
     .want = kIllegalIndex},
    {.name = "size two rotate one miss above",
     .nn = {3, 1},
     .k = 4,
     .want = kIllegalIndex},
    // size 3 - no rotation
    {.name = "size three no rotation miss below",
     .nn = {1, 3, 5},
     .k = 0,
     .want = kIllegalIndex},
    {.name = "size three no rotation miss between first and second",
     .nn = {1, 3, 5},
     .k = 2,
     .want = kIllegalIndex},
    {.name = "size three no rotation miss between second and third",
     .nn = {1, 3, 5},
     .k = 4,
     .want = kIllegalIndex},
    {.name = "size three no rotation miss above",
     .nn = {1, 3, 5},
     .k = 6,
     .want = kIllegalIndex},
    // size 3 - rotate one
    {.name = "size three rotate one miss below",
     .nn = {3, 5, 1},
     .k = 0,
     .want = kIllegalIndex},
    {.name = "size three rotate one miss between first and second",
     .nn = {3, 5, 1},
     .k = 2,
     .want = kIllegalIndex},
    {.name = "size three rotate one miss between second and third",
     .nn = {3, 5, 1},
     .k = 4,
     .want = kIllegalIndex},
    {.name = "size three rotate one miss above",
     .nn = {3, 5, 1},
     .k = 6,
     .want = kIllegalIndex},
    // size 3 - rotate two
    {.name = "size three rotate two miss below",
     .nn = {5, 1, 3},
     .k = 0,
     .want = kIllegalIndex},
    {.name = "size three rotate two miss between first and second",
     .nn = {5, 1, 3},
     .k = 2,
     .want = kIllegalIndex},
    {.name = "size three rotate two miss between second and third",
     .nn = {5, 1, 3},
     .k = 4,
     .want = kIllegalIndex},
    {.name = "size three rotate two miss above",
     .nn = {5, 1, 3},
     .k = 6,
     .want = kIllegalIndex},
};

INSTANTIATE_TEST_SUITE_P(FindTest, NegativeFindTest,
                         ValuesIn(kNegativeFindTestParams), NameTest);

}  // namespace
}  // namespace iq::rotated
