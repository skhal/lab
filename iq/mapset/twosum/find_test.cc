// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build nobuild

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "iq/mapset/twosum/find.h"

namespace iq::mapset::twosum {
namespace {

using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

struct FindTestParam {
  std::string name;
  std::vector<int> nn;
  int x;
  std::optional<Indices> want;

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    return os << "name: " << std::quoted(tp.name) << " nn: ["
              << absl::StrJoin(tp.nn, " ") << "]";
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const std::optional<Indices> got = Find(tp.nn, tp.x);

  EXPECT_THAT(got, Eq(tp.want));
}

const FindTestParam kFindTestParams[]{
    {.name = "empty", .nn = {}, .x = 4, .want = std::nullopt},
    {.name = "one item", .nn = {4}, .x = 4, .want = std::nullopt},
    {.name = "two items", .nn = {1, 3}, .x = 4, .want = std::make_pair(0, 1)},
    {.name = "two items not found", .nn = {1, 4}, .x = 4, .want = std::nullopt},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kFindTestParams),
                         [](const TestParamInfo<FindTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::mapset::twosum
