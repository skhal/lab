// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/19/midpoint/find.h"

#include <algorithm>
#include <cctype>
#include <initializer_list>
#include <memory>
#include <ostream>
#include <string>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::midpoint {
namespace {

using ::testing::Eq;
using ::testing::Field;
using ::testing::IsNull;
using ::testing::Pointee;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

std::vector<int> toVector(std::shared_ptr<Node> node) {
  std::vector<int> nn;
  for (; node; node = node->next) {
    nn.emplace_back(node->val);
  }
  return nn;
}

std::shared_ptr<Node> makeList(std::initializer_list<int> nn) {
  std::shared_ptr<Node> head;
  std::shared_ptr<Node> tail;
  for (int n : nn) {
    auto node = std::make_shared<Node>(n);
    if (!head) {
      head = node;
    } else {
      tail->next = node;
    }
    tail = node;
  }
  return head;
}

struct FindTestParam {
  std::string name;
  std::shared_ptr<Node> list;
  int want;

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    const std::vector<int> nn = toVector(tp.list);
    return os << "list: {" << absl::StrJoin(nn, ", ") << "}";
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

TEST(FindTest, Empty) {
  const std::shared_ptr<Node> got = Find(nullptr);

  EXPECT_THAT(got, IsNull());
}

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  const std::shared_ptr<Node> got = Find(tp.list);

  EXPECT_THAT(got, Pointee(Field(&Node::val, Eq(tp.want))));
}

const FindTestParam kFindTestParams[]{
    {.name = "one item", .list = makeList({1}), .want = 1},
    {.name = "two items", .list = makeList({1, 2}), .want = 2},
    {.name = "three items", .list = makeList({1, 2, 3}), .want = 2},
    {.name = "four items", .list = makeList({1, 2, 3, 4}), .want = 3},
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
}  // namespace iq::midpoint
