// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/2/threesum/solution.h"

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <iterator>
#include <ostream>
#include <sstream>
#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::threesum {
namespace {

using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

struct FindTestParam {
  std::string name;
  std::vector<int> nn;
  std::vector<Triplet> want;

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    std::ostringstream oss;
    std::copy(tp.nn.begin(), tp.nn.end(),
              std::ostream_iterator<int>(oss, ", "));
    return os << std::quoted(tp.name) << " nn: " << oss.str();
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

MATCHER(TupleEq, "") {
  const auto [a, b, c] = std::get<0>(arg);
  const auto [x, y, z] = std::get<1>(arg);
  std::vector<int> got{a, b, c};
  std::vector<int> want{x, y, z};
  ::testing::Matcher<std::vector<int>> matcher =
      ::testing::UnorderedElementsAreArray(want);
  return matcher.Matches(got);
}

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const std::vector<Triplet> got = Find(tp.nn);

  EXPECT_THAT(got, ::testing::UnorderedPointwise(TupleEq(), tp.want));
}

const FindTestParam kFindTestParams[]{
    // Negative tests
    {.name = "empty", .nn = {}, .want = {}},
    {.name = "one item", .nn = {1}, .want = {}},
    {.name = "two items", .nn = {1, 2}, .want = {}},
    {.name = "three positives", .nn = {1, 2, 3}, .want = {}},
    {.name = "three identical", .nn = {1, 1, 1}, .want = {}},
    {.name = "three negatives", .nn = {-1, -2, -3}, .want = {}},
    {.name = "three identical negatives", .nn = {-1, -1, -1}, .want = {}},
    // Positive tests
    // -- 3 items
    {.name = "three items", .nn = {1, 2, -3}, .want = {{-3, 1, 2}}},
    // -- 4 items
    {.name = "four items", .nn = {1, 2, -3, 4}, .want = {{-3, 1, 2}}},
    {.name = "four items duplicate low",
     .nn = {1, 2, -3, -3},
     .want = {{-3, 1, 2}}},
    {.name = "four items duplicate high",
     .nn = {1, 2, -3, 2},
     .want = {{-3, 1, 2}}},
    {.name = "four items duplicate middle",
     .nn = {1, 2, -3, 1},
     .want = {{-3, 1, 2}}},
    // -- 5 items
    {.name = "five items one triplet",
     .nn = {2, 4, -6, 3, 4},
     .want = {{-6, 2, 4}}},
    {
        .name = "five items two triplets",
        .nn = {2, 4, -6, 3, 3},
        .want = {{-6, 2, 4}, {-6, 3, 3}},
    }};

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
}  // namespace iq::threesum
