// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/reverse/reverse.h"

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

namespace iq::list::singly::reverse {
namespace {

using ::testing::ContainerEq;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

std::vector<int> makeVector(std::shared_ptr<Node> node) {
  std::vector<int> nn;
  while (node) {
    nn.push_back(node->Value);
    node = node->Next;
  }
  return nn;
}

struct ReverseTestParam {
  std::string name;
  std::shared_ptr<Node> list;
  std::shared_ptr<Node> want;

  friend std::ostream& operator<<(std::ostream& os,
                                  const ReverseTestParam& tp) {
    const std::vector<int> nn = makeVector(tp.list);
    os << "list: {" << absl::StrJoin(nn, ", ") << "}";
    return os;
  }
};

class ReverseTest : public TestWithParam<ReverseTestParam> {};

std::shared_ptr<Node> makeList(std::initializer_list<int> nn) {
  std::shared_ptr<Node> head = nullptr;
  std::shared_ptr<Node> tail = nullptr;
  for (int n : nn) {
    std::shared_ptr<Node> node = std::make_shared<Node>(n, nullptr);
    if (head == nullptr) {
      head = node;
      tail = node;
      continue;
    }
    tail->Next = node;
    tail = node;
  }
  return head;
}

MATCHER_P(ListEq, list, "") {
  const std::vector<int> got = makeVector(arg);
  const std::vector<int> want = makeVector(list);
  EXPECT_THAT(got, ContainerEq(want));
  return true;
}

TEST_P(ReverseTest, Test) {
  const ReverseTestParam& tp = GetParam();

  const std::shared_ptr<Node> got = Reverse(tp.list);

  EXPECT_THAT(got, ListEq(tp.want));
}

const ReverseTestParam kReverseTestParams[]{
    {.name = "empty", .list = nullptr, .want = nullptr},
    {.name = "one item", .list = makeList({1}), .want = makeList({1})},
    {.name = "two items", .list = makeList({1, 2}), .want = makeList({2, 1})},
    {.name = "three items",
     .list = makeList({1, 2, 3}),
     .want = makeList({3, 2, 1})},
};

INSTANTIATE_TEST_SUITE_P(
    ReverseTest, ReverseTest, ValuesIn(kReverseTestParams),
    [](const testing::TestParamInfo<ReverseTestParam>& info) {
      std::string name = info.param.name;
      std::replace_if(
          name.begin(), name.end(), [](char c) { return !std::isalnum(c); },
          kCharUnderscore);
      return name;
    });

}  // namespace
}  // namespace iq::list::singly::reverse
