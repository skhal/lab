// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/25/bounds/find.h"

#include <algorithm>
#include <cctype>
#include <ostream>
#include <string>
#include <string_view>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::bounds {
namespace {

using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct FindTestParam {
  std::string_view name;
  std::vector<int> nn;
  int n;
  Bounds want;
};

std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
  return os << "{ nn:[" << absl::StrJoin(tp.nn, " ") << "] n:" << tp.n << " }";
}

class FindTest : public TestWithParam<FindTestParam> {};

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const Bounds got = Find(tp.nn, tp.n);

  EXPECT_THAT(got, Eq(tp.want));
}

const FindTestParam kFindTestParams[]{
    // size 0
    {.name = "size 0 empty", .nn = {}, .n = 1, .want = kIllegalBounds},
    // size 1
    {.name = "size 1 hit", .nn = {1}, .n = 1, .want = {0, 0}},
    {.name = "size 1 miss before", .nn = {1}, .n = 0, .want = kIllegalBounds},
    {.name = "size 1 miss after", .nn = {1}, .n = 2, .want = kIllegalBounds},
    // size 2
    {.name = "size 2 no dups hit first", .nn = {1, 3}, .n = 1, .want = {0, 0}},
    {.name = "size 2 no dups hit second", .nn = {1, 3}, .n = 3, .want = {1, 1}},
    {.name = "size 2 no dups miss before",
     .nn = {1, 3},
     .n = 0,
     .want = kIllegalBounds},
    {.name = "size 2 no dups miss between",
     .nn = {1, 3},
     .n = 2,
     .want = kIllegalBounds},
    {.name = "size 2 no dups miss after",
     .nn = {1, 3},
     .n = 4,
     .want = kIllegalBounds},
    {.name = "size 2 dups hit", .nn = {1, 1}, .n = 1, .want = {0, 1}},
    {.name = "size 2 dups miss before",
     .nn = {1, 1},
     .n = 0,
     .want = kIllegalBounds},
    {.name = "size 2 dups miss after",
     .nn = {1, 1},
     .n = 2,
     .want = kIllegalBounds},
    // size 3
    {.name = "size 3 no dups hit first",
     .nn = {1, 3, 5},
     .n = 1,
     .want = {0, 0}},
    {.name = "size 3 no dups hit second",
     .nn = {1, 3, 5},
     .n = 3,
     .want = {1, 1}},
    {.name = "size 3 no dups hit third",
     .nn = {1, 3, 5},
     .n = 5,
     .want = {2, 2}},
    {.name = "size 3 no dups miss before",
     .nn = {1, 3, 5},
     .n = 0,
     .want = kIllegalBounds},
    {.name = "size 3 no dups miss between first and second",
     .nn = {1, 3, 5},
     .n = 2,
     .want = kIllegalBounds},
    {.name = "size 3 no dups miss between second and third",
     .nn = {1, 3, 5},
     .n = 4,
     .want = kIllegalBounds},
    {.name = "size 3 no dups miss after",
     .nn = {1, 3, 5},
     .n = 6,
     .want = kIllegalBounds},
    {.name = "size 3 first two dups hit first",
     .nn = {1, 1, 3},
     .n = 1,
     .want = {0, 1}},
    {.name = "size 3 first two dups hit third",
     .nn = {1, 1, 3},
     .n = 3,
     .want = {2, 2}},
    {.name = "size 3 first two dups miss before",
     .nn = {1, 1, 3},
     .n = 0,
     .want = kIllegalBounds},
    {.name = "size 3 first two dups miss between second and third",
     .nn = {1, 1, 3},
     .n = 2,
     .want = kIllegalBounds},
    {.name = "size 3 first two dups miss after",
     .nn = {1, 1, 3},
     .n = 4,
     .want = kIllegalBounds},
    {.name = "size 3 last two dups hit first",
     .nn = {1, 3, 3},
     .n = 1,
     .want = {0, 0}},
    {.name = "size 3 last two dups hit third",
     .nn = {1, 3, 3},
     .n = 3,
     .want = {1, 2}},
    {.name = "size 3 last two dups miss before",
     .nn = {1, 3, 3},
     .n = 0,
     .want = kIllegalBounds},
    {.name = "size 3 last two dups miss between first and second ",
     .nn = {1, 3, 3},
     .n = 2,
     .want = kIllegalBounds},
    {.name = "size 3 last two dups miss after",
     .nn = {1, 3, 3},
     .n = 4,
     .want = kIllegalBounds},

    {.name = "size 3 all dups hit", .nn = {1, 1, 1}, .n = 1, .want = {0, 2}},
    {.name = "size 3 all dups miss before",
     .nn = {1, 1, 1},
     .n = 0,
     .want = kIllegalBounds},
    {.name = "size 3 all dups miss after",
     .nn = {1, 1, 1},
     .n = 2,
     .want = kIllegalBounds},
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
}  // namespace iq::bounds
