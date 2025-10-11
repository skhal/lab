// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build nobuild

#include <algorithm>
#include <cctype>
#include <iomanip>
#include <string>
#include <string_view>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "iq/mapset/chain/find.h"

namespace iq::mapset::chain {
namespace {

using ::testing::AnyOfArray;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUdnerscore = '_';
constexpr std::string_view kStrSpace = " ";

using Chain = std::vector<int>;

struct FindTestParam {
  std::string name;
  std::vector<int> nn;
  std::vector<Chain> want;

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    return os << " name: " << std::quoted(tp.name) << " nn: ["
              << absl::StrJoin(tp.nn, kStrSpace) << "]";
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const std::vector<int> got = Find(tp.nn);

  EXPECT_THAT(got, AnyOfArray(tp.want));
}

template <typename... Args>
std::vector<int> makeChain(Args&&... args) {
  return {args...};
}

const FindTestParam kFindTestParams[]{
    {.name = "empty", .nn = {}, .want = {makeChain()}},
    {.name = "one item", .nn = {1}, .want = {makeChain(1)}},
    {.name = "two items one chain", .nn = {1, 2}, .want = {makeChain(1, 2)}},
    {.name = "two items one chain reversed",
     .nn = {2, 1},
     .want = {makeChain(1, 2)}},
    {.name = "two items two chains",
     .nn = {1, 3},
     .want = {makeChain(1), makeChain(3)}},
    {.name = "two chains same size",
     .nn = {1, 5, 2, 4},
     .want = {makeChain(1, 2), makeChain(4, 5)}},
    {.name = "two chains different size",
     .nn = {1, 7, 5, 2, 6},
     .want = {makeChain(5, 6, 7)}},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kFindTestParams),
                         [](const TestParamInfo<FindTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUdnerscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::mapset::chain
