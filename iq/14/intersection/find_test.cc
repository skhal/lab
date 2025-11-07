// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/14/intersection/find.h"

#include <initializer_list>
#include <iterator>
#include <memory>
#include <ostream>
#include <string>
#include <string_view>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::intersection {
namespace {

using ::testing::ContainerEq;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr std::string_view kSepList = ", ";

std::vector<int> toVector([[maybe_unused]] std::shared_ptr<Node> list) {
  std::vector<int> nn;
  for (; list; list = list->next) {
    nn.emplace_back(list->val);
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
      tail = node;
      continue;
    }
    tail->next = node;
    tail = node;
  }
  return head;
}

std::shared_ptr<Node> makeList(std::initializer_list<int> nn,
                               std::shared_ptr<Node> tail) {
  if (std::empty(nn)) {
    return tail;
  }
  std::shared_ptr<Node> head = makeList(nn);
  std::shared_ptr<Node> node = head;
  while (node->next) {
    node = node->next;
  }
  node->next = tail;
  return head;
}

MATCHER_P(ListEq, want, "") {
  const std::vector<int> got = toVector(arg);
  EXPECT_THAT(got, ContainerEq(want));
  return true;
}

struct FindTestParam {
  std::string name;
  std::shared_ptr<Node> a;
  std::shared_ptr<Node> b;
  std::vector<int> want;

  friend std::ostream& operator<<(std::ostream& os, const FindTestParam& tp) {
    return os << "a: {" << absl::StrJoin(toVector(tp.a), kSepList) << "} b: {"
              << absl::StrJoin(toVector(tp.b), kSepList) << "}";
  }
};

class FindTest : public TestWithParam<FindTestParam> {};

TEST_P(FindTest, Test) {
  const FindTestParam& tp = GetParam();

  std::shared_ptr<Node> got = Find(tp.a, tp.b);

  EXPECT_THAT(got, ListEq(tp.want));
}

const struct {
  std::shared_ptr<Node> one;
  std::shared_ptr<Node> two;
} kSharedList{
    .one = makeList({1}),
    .two = makeList({1, 2}),
};

const FindTestParam kFindTestParams[]{
    {.name = "empty", .a = nullptr, .b = nullptr, .want = {}},
    {.name = "empty and one item",
     .a = nullptr,
     .b = kSharedList.one,
     .want = {}},
    {.name = "one item and empty",
     .a = kSharedList.one,
     .b = nullptr,
     .want = {}},
    // --
    {.name = "one item no intersection",
     .a = makeList({1}),
     .b = makeList({1}),
     .want = {}},
    {.name = "one item intersection all",
     .a = kSharedList.one,
     .b = kSharedList.one,
     .want = toVector(kSharedList.one)},
    // --
    {.name = "two items no intersection",
     .a = makeList({10, 1}),
     .b = makeList({20, 1}),
     .want = {}},
    {.name = "two items intersection one item",
     .a = makeList({10}, kSharedList.one),
     .b = makeList({20}, kSharedList.one),
     .want = toVector(kSharedList.one)},
    {.name = "two same items intersection one item",
     .a = makeList({2}, kSharedList.one),
     .b = makeList({2}, kSharedList.one),
     .want = toVector(kSharedList.one)},
    {.name = "two items intersection all",
     .a = kSharedList.two,
     .b = kSharedList.two,
     .want = toVector(kSharedList.two)},
    {.name = "two items and one item intersection one item",
     .a = makeList({10}, kSharedList.one),
     .b = kSharedList.one,
     .want = toVector(kSharedList.one)},
    {.name = "one item and two items intersection one item",
     .a = kSharedList.one,
     .b = makeList({20}, kSharedList.one),
     .want = toVector(kSharedList.one)},
};

INSTANTIATE_TEST_SUITE_P(FindTest, FindTest, ValuesIn(kFindTestParams));

}  // namespace
}  // namespace iq::intersection
