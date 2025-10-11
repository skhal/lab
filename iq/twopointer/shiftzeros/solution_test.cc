// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build ignore

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
#include "iq/twopointer/shiftzeros/solution.h"

namespace iq::twopointer::shiftzeros {
namespace {

using ::testing::ElementsAreArray;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

struct ShiftTestParam {
  std::string name;
  std::vector<int> nn;
  std::vector<int> want;

  friend std::ostream& operator<<(std::ostream& os, const ShiftTestParam& tp) {
    std::ostringstream oss;
    std::copy(tp.nn.begin(), tp.nn.end(), std::ostream_iterator<int>(oss, ","));
    return os << std::quoted(tp.name) << " nn: " << oss.str();
  }
};

class ShiftTest : public TestWithParam<ShiftTestParam> {};

TEST_P(ShiftTest, Test) {
  const ShiftTestParam& tp = GetParam();

  std::vector<int> got = tp.nn;
  Shift(got);

  EXPECT_THAT(got, ElementsAreArray(tp.want));
}

const ShiftTestParam kShiftTestParams[]{
    {.name = "empty", .nn = {}, .want = {}},
    // 1 item
    {.name = "one item not zero", .nn = {1}, .want = {1}},
    {.name = "one item zero", .nn = {0}, .want = {0}},
    // 2 items
    {.name = "two items not zero", .nn = {1, 2}, .want = {1, 2}},
    {.name = "two items first zero", .nn = {0, 1}, .want = {1, 0}},
    {.name = "two items second zero", .nn = {1, 0}, .want = {1, 0}},
    {.name = "two zeros", .nn = {0, 0}, .want = {0, 0}},
    // 3 itmes
    {.name = "three items not zeros", .nn = {1, 2, 3}, .want = {1, 2, 3}},
    {.name = "three items first zero", .nn = {0, 2, 3}, .want = {2, 3, 0}},
    {.name = "three items second zero", .nn = {1, 0, 3}, .want = {1, 3, 0}},
    {.name = "three items third zero", .nn = {1, 2, 0}, .want = {1, 2, 0}},
    {.name = "three items first and second zero",
     .nn = {0, 0, 3},
     .want = {3, 0, 0}},
    {.name = "three items first and third zero",
     .nn = {0, 2, 0},
     .want = {2, 0, 0}},
    {.name = "three items second and third zero",
     .nn = {1, 0, 0},
     .want = {1, 0, 0}},
    {.name = "three zeros", .nn = {0, 0, 0}, .want = {0, 0, 0}},
};

INSTANTIATE_TEST_SUITE_P(ShiftTest, ShiftTest, ValuesIn(kShiftTestParams),
                         [](const TestParamInfo<ShiftTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::twopointer::shiftzeros
