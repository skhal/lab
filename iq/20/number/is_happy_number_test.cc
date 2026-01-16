// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/20/number/is_happy_number.h"

#include <ostream>

#include "absl/strings/str_format.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::number {
namespace {

using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

struct IsHappyNumberTestParam {
  int n;
  bool want;

  friend std::ostream& operator<<(std::ostream& os,
                                  const IsHappyNumberTestParam& tp) {
    return os << " n: " << tp.n;
  }
};

class IsHappyNumberTest : public TestWithParam<IsHappyNumberTestParam> {};

TEST_P(IsHappyNumberTest, Test) {
  const IsHappyNumberTestParam& tp = GetParam();

  const bool got = IsHappyNumber(tp.n);

  EXPECT_THAT(got, Eq(tp.want));
}

const IsHappyNumberTestParam kIsHappyNumberTestParam[]{
    // Happy numbers
    {.n = 1, .want = true},
    {.n = 7, .want = true},
    {.n = 10, .want = true},
    {.n = 208, .want = true},
    {.n = 931, .want = true},
    // Not happy numbers
    {.n = 2, .want = false},
    {.n = 168, .want = false},
    {.n = 936, .want = false},
};

INSTANTIATE_TEST_SUITE_P(IsHappyNumberTest, IsHappyNumberTest,
                         ValuesIn(kIsHappyNumberTestParam),
                         [](const TestParamInfo<IsHappyNumberTestParam>& info) {
                           return absl::StrFormat("%d", info.param.n);
                         });

}  // namespace
}  // namespace iq::number
