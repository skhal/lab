// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/24/insert/find.h"

#include <algorithm>
#include <cctype>
#include <cstddef>
#include <ostream>
#include <string>
#include <string_view>
#include <vector>

#include "absl/strings/str_join.h"
#include "gtest/gtest.h"

namespace iq::insert {
namespace {

using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindInsertIndexTestParam {
  std::string_view name;
  std::vector<int> nn;
  int n;
  std::size_t want;
};

std::ostream& operator<<(std::ostream& os, const FindInsertIndexTestParam& tp) {
  return os << "{ nn:[" << absl::StrJoin(tp.nn, " ") << "] n:" << tp.n << " }";
}

class FindInsertIndexTest : public TestWithParam<FindInsertIndexTestParam> {};

TEST(FindInsertIndexTest, Example) {
  const std::size_t got = FindInsertIndex({1, 2, 3}, 2);

  EXPECT_EQ(got, static_cast<std::size_t>(1));
}

TEST_P(FindInsertIndexTest, Test) {
  const FindInsertIndexTestParam& tp = GetParam();

  const std::size_t got = FindInsertIndex(tp.nn, tp.n);

  EXPECT_EQ(got, tp.want);
}

const FindInsertIndexTestParam kFindInsertIndexTestParam[]{
    {.name = "empty", .nn = {}, .n = 1, .want = 0},
    // size 1
    {.name = "one item hit", .nn = {1}, .n = 1, .want = 0},
    {.name = "one item insert before", .nn = {1}, .n = 0, .want = 0},
    {.name = "one item insert after", .nn = {1}, .n = 2, .want = 1},
    // size 2
    {.name = "two items hit first", .nn = {1, 3}, .n = 1, .want = 0},
    {.name = "two items hit second", .nn = {1, 3}, .n = 3, .want = 1},
    {.name = "two items insert before", .nn = {1, 3}, .n = 0, .want = 0},
    {.name = "two items insert middle", .nn = {1, 3}, .n = 2, .want = 1},
    {.name = "two items insert after", .nn = {1, 3}, .n = 4, .want = 2},
    // size 3
    {.name = "three items hit first", .nn = {1, 3, 5}, .n = 1, .want = 0},
    {.name = "three items hit second", .nn = {1, 3, 5}, .n = 3, .want = 1},
    {.name = "three items hit third", .nn = {1, 3, 5}, .n = 5, .want = 2},
    {.name = "three items insert before", .nn = {1, 3, 5}, .n = 0, .want = 0},
    {.name = "three items insert between first and second",
     .nn = {1, 3, 5},
     .n = 2,
     .want = 1},
    {.name = "three items insert between second and third",
     .nn = {1, 3, 5},
     .n = 4,
     .want = 2},
    {.name = "three items insert after", .nn = {1, 3, 5}, .n = 6, .want = 3},
    // size 4
    {.name = "four items hit first", .nn = {1, 3, 5, 7}, .n = 1, .want = 0},
    {.name = "four items hit second", .nn = {1, 3, 5, 7}, .n = 3, .want = 1},
    {.name = "four items hit third", .nn = {1, 3, 5, 7}, .n = 5, .want = 2},
    {.name = "four items hit fourth", .nn = {1, 3, 5, 7}, .n = 7, .want = 3},
    {.name = "four items insert before", .nn = {1, 3, 5, 7}, .n = 0, .want = 0},
    {.name = "four items insert between first and second",
     .nn = {1, 3, 5, 7},
     .n = 2,
     .want = 1},
    {.name = "four items insert between second and third",
     .nn = {1, 3, 5, 7},
     .n = 4,
     .want = 2},
    {.name = "four items insert between third and fourth",
     .nn = {1, 3, 5, 7},
     .n = 6,
     .want = 3},
    {.name = "four items insert after", .nn = {1, 3, 5, 7}, .n = 8, .want = 4},
};

INSTANTIATE_TEST_SUITE_P(
    FindInsertIndexTest, FindInsertIndexTest,
    ValuesIn(kFindInsertIndexTestParam),
    [](const TestParamInfo<FindInsertIndexTestParam>& info) {
      std::string name(info.param.name);
      std::replace_if(
          name.begin(), name.end(), [](char c) { return !std::isalnum(c); },
          '_');
      return name;
    });

}  // namespace
}  // namespace iq::insert
