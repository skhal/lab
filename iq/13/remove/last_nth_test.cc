// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/13/remove/last_nth.h"

#include <algorithm>
#include <cctype>
#include <cstddef>
#include <initializer_list>
#include <memory>
#include <ostream>
#include <string>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::remove {
namespace {

using ::testing::ContainerEq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

std::shared_ptr<Node> makeList(std::initializer_list<int> nn) {
  std::shared_ptr<Node> head;
  std::shared_ptr<Node> tail;
  for (int n : nn) {
    std::shared_ptr<Node> node = std::make_shared<Node>(n);
    if (head == nullptr) {
      head = node;
      tail = node;
      continue;
    }
    tail->next = node;
    tail = std::move(node);
  }
  return head;
}

std::vector<int> makeVector(std::shared_ptr<Node> list) {
  std::vector<int> nn;
  for (; list; list = list->next) {
    nn.emplace_back(list->value);
  }
  return nn;
}

struct LastNthTestParam {
  std::string name;
  std::shared_ptr<Node> list;
  std::size_t n;
  std::shared_ptr<Node> want;

  friend std::ostream& operator<<(std::ostream& os,
                                  const LastNthTestParam& tp) {
    const std::vector<int> nn = makeVector(tp.list);
    return os << "list: {" << absl::StrJoin(nn, ", ") << "} n: " << tp.n;
  }
};

class LastNthTest : public TestWithParam<LastNthTestParam> {};

MATCHER_P(ListEq, expected, "") {
  const std::vector<int> got = makeVector(arg);
  const std::vector<int> want = makeVector(expected);
  EXPECT_THAT(got, ContainerEq(want));
  return true;
}

TEST_P(LastNthTest, Test) {
  const LastNthTestParam& tp = GetParam();

  std::shared_ptr<Node> got = RemoveLastNth(tp.list, tp.n);

  EXPECT_THAT(got, ListEq(tp.want));
}

const LastNthTestParam kLastNthTestParams[]{
    {.name = "empty", .list = nullptr, .n = 1, .want = nullptr},
    // size 1
    {.name = "one item remove head",
     .list = makeList({1}),
     .n = 1,
     .want = nullptr},
    {.name = "one item insufficient items",
     .list = makeList({1}),
     .n = 2,
     .want = makeList({1})},
    // size 2
    {.name = "two items remove head",
     .list = makeList({1, 2}),
     .n = 2,
     .want = makeList({2})},
    {.name = "two items remove tail",
     .list = makeList({1, 2}),
     .n = 1,
     .want = makeList({1})},
    {.name = "two items insufficient items",
     .list = makeList({1, 2}),
     .n = 3,
     .want = makeList({1, 2})},
    // size 3
    {.name = "three items remove head",
     .list = makeList({1, 2, 3}),
     .n = 3,
     .want = makeList({2, 3})},
    {.name = "three items remove tail",
     .list = makeList({1, 2, 3}),
     .n = 1,
     .want = makeList({1, 2})},
    {.name = "three items remove first",
     .list = makeList({1, 2, 3}),
     .n = 2,
     .want = makeList({1, 3})},
    {.name = "three items insufficient items",
     .list = makeList({1, 2, 3}),
     .n = 4,
     .want = makeList({1, 2, 3})},
};

INSTANTIATE_TEST_SUITE_P(LastNthTest, LastNthTest, ValuesIn(kLastNthTestParams),
                         [](const TestParamInfo<LastNthTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::remove
