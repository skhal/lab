// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/16/palindrome/is.h"

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

namespace iq::palindrome {
namespace {

using ::testing::Eq;
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
    tail = node;
  }
  return head;
}

std::vector<int> makeVector(std::shared_ptr<Node> list) {
  std::vector<int> nn;
  for (; list; list = list->next) {
    nn.emplace_back(list->val);
  }
  return nn;
}

struct IsTestParam {
  std::string name;
  std::shared_ptr<Node> list;
  bool want;

  friend std::ostream& operator<<(std::ostream& os, const IsTestParam& tp) {
    const std::vector<int> nn = makeVector(tp.list);
    return os << "list: {" << absl::StrJoin(nn, ", ") << "}";
  }
};

class IsTest : public TestWithParam<IsTestParam> {};

TEST_P(IsTest, Test) {
  const IsTestParam& tp = GetParam();

  const bool got = Is(tp.list);

  EXPECT_THAT(got, Eq(tp.want));
}

const IsTestParam kIsTestParams[]{
    {.name = "empty", .list = makeList({}), .want = true},
    {.name = "one item", .list = makeList({1}), .want = true},
    {.name = "two items palindrome", .list = makeList({1, 1}), .want = true},
    {.name = "two items not palindrome",
     .list = makeList({1, 2}),
     .want = false},
    {.name = "three items palindrome",
     .list = makeList({1, 2, 1}),
     .want = true},
    {.name = "three items not palindrome",
     .list = makeList({1, 2, 3}),
     .want = false},
};

INSTANTIATE_TEST_SUITE_P(IsTest, IsTest, ValuesIn(kIsTestParams),
                         [](const TestParamInfo<IsTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::palindrome
