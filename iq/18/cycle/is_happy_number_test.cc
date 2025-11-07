// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/18/cycle/is_happy_number.h"

#include <ostream>

#include "absl/strings/str_format.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::cycle {
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
    {.n = 188, .want = true},
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
}  // namespace iq::cycle
