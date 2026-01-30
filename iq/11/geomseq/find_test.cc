// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/11/geomseq/find.h"

#include <cctype>
#include <iomanip>
#include <ostream>
#include <sstream>
#include <string>
#include <string_view>
#include <tuple>
#include <vector>

#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::geomseq {

// Must be in the same namespace with Triplet:
// https://google.github.io/googletest/advanced.html#teaching-googletest-how-to-print-your-values
template <typename Sink>
void AbslStringify(Sink& sink, const Triplet& t) {
  std::ostringstream os;
  os << t;
  absl::Format(&sink, os.str());
}

namespace {

using ::testing::Pointwise;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr std::string_view kStrSpace = " ";
constexpr char kCharUnderscore = '_';

struct FindTestParam {
  std::string name;
  std::vector<int> nn = {};
  int ratio;
  std::vector<Triplet> want = {};

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    return os << "name: " << std::quoted(tp.name) << " nn: ["
              << absl::StrJoin(tp.nn, kStrSpace) << "]";
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

MATCHER(TripletEq, "equal triplets") {
  const Triplet& got = std::get<0>(arg);
  const Triplet& want = std::get<1>(arg);
  return std::tie(got.i, got.j, got.k) == std::tie(want.i, want.j, want.k);
}

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const std::vector<Triplet> got = Find(tp.nn, tp.ratio);

  EXPECT_THAT(got, Pointwise(TripletEq(), tp.want));
}

const FindTestParam kFindTestParams[]{
    {.name = "empty", .ratio = 2},
    {.name = "one item", .nn = {1}, .ratio = 2},
    {.name = "two items", .nn = {1, 2}, .ratio = 2},
    {.name = "three items no seq", .nn = {1, 2, 3}, .ratio = 2},
    {
        .name = "three items seq",
        .nn = {1, 2, 4},
        .ratio = 2,
        .want =
            {
                {0, 1, 2},
            },
    },
    {
        .name = "four items two seq",
        .nn = {1, 2, 2, 4},
        .ratio = 2,
        .want =
            {
                {0, 1, 3},
                {0, 2, 3},
            },
    },
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
}  // namespace iq::geomseq
