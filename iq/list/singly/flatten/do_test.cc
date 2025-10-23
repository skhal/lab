// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/flatten/do.h"

#include <cctype>
#include <initializer_list>
#include <memory>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::list::singly::flatten {
namespace {

using ::testing::ContainerEq;
using ::testing::IsEmpty;

std::shared_ptr<TreeNode> makeTree(std::initializer_list<int> nn) {
  std::shared_ptr<TreeNode> head;
  std::shared_ptr<TreeNode> tail;
  for (int n : nn) {
    auto node = std::make_shared<TreeNode>(n);
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

std::shared_ptr<TreeNode> getNode(std::shared_ptr<TreeNode> node, int val) {
  for (; node; node = node->next) {
    if (node->val == val) {
      return node;
    }
  }
  return nullptr;
}

std::vector<int> toVector(std::shared_ptr<Node> list) {
  std::vector<int> nn;
  for (; list; list = list->next) {
    nn.emplace_back(list->val);
  }
  return nn;
}

MATCHER_P(ListEq, want, "") {
  const std::vector<int> got = toVector(arg);
  EXPECT_THAT(got, ContainerEq(want));
  return true;
}

MATCHER(IsEmptyList, "") {
  const std::vector<int> got = toVector(arg);
  EXPECT_THAT(got, IsEmpty());
  return true;
}

TEST(DoTest, Empty) {
  auto tree = makeTree({});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, IsEmptyList());
}

TEST(DoTest, FlatOneItem) {
  auto tree = makeTree({1});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1}));
}

TEST(DoTest, FlatTwoItems) {
  auto tree = makeTree({1, 2});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2}));
}

TEST(DoTest, FlatThreeItems) {
  auto tree = makeTree({1, 2, 3});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2, 3}));
}

TEST(DoTest, Level1OneLevel2One) {
  auto tree = makeTree({1});
  getNode(tree, 1)->child = makeTree({2});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2}));
}

TEST(DoTest, Level1TwoLeve2OneOnFirst) {
  auto tree = makeTree({1, 2});
  getNode(tree, 1)->child = makeTree({3});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2, 3}));
}

TEST(DoTest, Level1TwoLeve2OneOnSecond) {
  auto tree = makeTree({1, 2});
  getNode(tree, 2)->child = makeTree({3});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2, 3}));
}

TEST(DoTest, Level1TwoLeve2TwoOnFirst) {
  auto tree = makeTree({1, 2});
  getNode(tree, 1)->child = makeTree({3, 4});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2, 3, 4}));
}

TEST(DoTest, Level1TwoLeve2TwoOnSecond) {
  auto tree = makeTree({1, 2});
  getNode(tree, 2)->child = makeTree({3, 4});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2, 3, 4}));
}

TEST(DoTest, Level1TwoLeve2TwoOnEach) {
  auto tree = makeTree({1, 2});
  getNode(tree, 1)->child = makeTree({3});
  getNode(tree, 2)->child = makeTree({4});

  std::shared_ptr<Node> got = Do(tree);

  EXPECT_THAT(got, ListEq(std::vector<int>{1, 2, 3, 4}));
}

}  // namespace
}  // namespace iq::list::singly::flatten
