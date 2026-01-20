// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/27/lower/find.h"

#include <ostream>
#include <string_view>
#include <vector>

#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::lower {
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
  return os << absl::StrFormat("{ nn:[%s] k:%d }", absl::StrJoin(tp.nn, " "),
                               tp.k);
}

class FindTest : public TestWithParam<FindTestParam> {};

using PositiveFindTest = FindTest;

TEST_P(PositiveFindTest, Test) {
  const FindTestParam& tp = GetParam();

  const absl::StatusOr<int> got = Find(tp.nn, tp.k);

  EXPECT_THAT(got, IsOkAndHolds(tp.want));
}

using NegativeFindTest = FindTest;

TEST_P(NegativeFindTest, Test) {
  const FindTestParam& tp = GetParam();

  const absl::StatusOr<int> got = Find(tp.nn, tp.k);

  EXPECT_THAT(got, Not(IsOk()));
}

const FindTestParam kPositiveFindTestParams[]{
    // size 1
    {.name = "size one below", .nn = {1}, .k = 0, .want = 1},
    {.name = "size one match", .nn = {1}, .k = 1, .want = 1},
    // size 2
    {.name = "size two below", .nn = {1, 3}, .k = 0, .want = 1},
    {.name = "size two match first", .nn = {1, 3}, .k = 1, .want = 1},
    {.name = "size two below second", .nn = {1, 3}, .k = 2, .want = 3},
    {.name = "size two match second", .nn = {1, 3}, .k = 3, .want = 3},
    // size 3
    {.name = "size three below", .nn = {1, 3, 5}, .k = 0, .want = 1},
    {.name = "size three match first", .nn = {1, 3, 5}, .k = 1, .want = 1},
    {.name = "size three below second", .nn = {1, 3, 5}, .k = 2, .want = 3},
    {.name = "size three match second", .nn = {1, 3, 5}, .k = 3, .want = 3},
    {.name = "size three below third", .nn = {1, 3, 5}, .k = 4, .want = 5},
    {.name = "size three match third", .nn = {1, 3, 5}, .k = 5, .want = 5},
};

INSTANTIATE_TEST_SUITE_P(FindTest, PositiveFindTest,
                         ValuesIn(kPositiveFindTestParams));

inline constexpr int kIllegalValue = 0;  // unused in test - can be any value

const FindTestParam kNegativeFindTestParams[]{
    {.name = "size one above", .nn = {1}, .k = 2, .want = kIllegalValue},
    {.name = "size two above", .nn = {1, 3}, .k = 4, .want = kIllegalValue},
    {.name = "size three above",
     .nn = {1, 3, 5},
     .k = 6,
     .want = kIllegalValue},
};

INSTANTIATE_TEST_SUITE_P(FindTest, NegativeFindTest,
                         ValuesIn(kNegativeFindTestParams));

}  // namespace
}  // namespace iq::lower
