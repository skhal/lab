// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build nobuild

#include <algorithm>
#include <string>
#include <vector>

#include "gtest/gtest.h"
#include "iq/twopointer/container/solution.h"

namespace iq::twopointer::container {
namespace {

using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharSpace = ' ';
constexpr char kCharDash = '-';
constexpr char kCharUnderscore = '_';

struct FindTestParam {
  std::string name;
  std::vector<int> nn;
  Volume want;
};

class FindTest : public TestWithParam<FindTestParam> {};

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const Volume got = Find(tp.nn);

  ASSERT_EQ(got, tp.want);
}

const std::vector<FindTestParam> kTestParams = {
    {.name = "empty", .nn = {}, .want = Volume{0}},
    {.name = "one item", .nn = {1}, .want = Volume{0}},
    {.name = "two items same non-zero", .nn = {1, 1}, .want = Volume{1}},
    {.name = "two items non-zero different", .nn = {1, 2}, .want = Volume{1}},
    {.name = "three items flat", .nn = {1, 1, 1}, .want = Volume{2}},
    {.name = "three items ascending", .nn = {1, 2, 3}, .want = Volume{2}},
    {.name = "three items descending", .nn = {3, 2, 1}, .want = Volume{2}},
    {.name = "four items same edges", .nn = {1, 4, 3, 1}, .want = Volume{3}},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kTestParams),
                         [](const TestParamInfo<FindTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](const char c) {
                                 return c == kCharSpace || c == kCharDash;
                               },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::twopointer::container
