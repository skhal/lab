// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build nobuild

#include <algorithm>
#include <optional>
#include <utility>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "iq/twopointer/twosum/solution.h"

namespace iq::twopointer::twosum {
namespace {

using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::ValuesIn;

struct FindTestParam {
  std::string name;
  std::vector<int> nn;
  int x;
  std::optional<std::pair<int, int>> want;
};

class FindTest : public ::testing::TestWithParam<FindTestParam> {};

TEST_P(FindTest, Demo) {
  const FindTestParam tc = GetParam();

  const std::optional<std::pair<int, int>> got = Find(tc.nn, tc.x);

  EXPECT_THAT(got, Eq(tc.want));
}

constexpr char kCharSpace = ' ';
constexpr char kCharUnderscore = '_';

const std::vector<FindTestParam> kFindTestCases{
    {.name = "empty input", .nn = {}, .x = 1, .want = std::nullopt},
    {.name = "one element", .nn = {1}, .x = 1, .want = std::nullopt},
    {.name = "no match", .nn = {1, 2}, .x = 4, .want = std::nullopt},
    {.name = "match", .nn = {1, 2}, .x = 3, .want = std::make_pair(0, 1)},
    {.name = "first match",
     .nn = {1, 1, 2},
     .x = 3,
     .want = std::make_pair(0, 2)},
    {.name = "move first index",
     .nn = {1, 2, 3, 4},
     .x = 6,
     .want = std::make_pair(1, 3)},
    {.name = "move second index",
     .nn = {1, 2, 3, 4},
     .x = 4,
     .want = std::make_pair(0, 2)},
    {.name = "negative values",
     .nn = {-1, 2, 3, 4},
     .x = 2,
     .want = std::make_pair(0, 2)},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kFindTestCases),
                         [](const TestParamInfo<FindTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](const char c) { return c == kCharSpace; },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::twopointer::twosum
