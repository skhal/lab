// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/cycle/is.h"

#include <memory>

#include "gtest/gtest.h"

namespace iq::list::singly::cycle {
namespace {

TEST(IsTest, empty) {
  std::shared_ptr<Node> list;

  const bool got = Is(list);

  EXPECT_FALSE(got);
}

TEST(IsTest, oneItemNoCycle) {
  std::shared_ptr<Node> list = std::make_shared<Node>(1);

  const bool got = Is(list);

  EXPECT_FALSE(got);
}

TEST(IsTest, oneItemCycle) {
  std::shared_ptr<Node> list = std::make_shared<Node>(1);
  list->next = list;

  const bool got = Is(list);

  EXPECT_TRUE(got);
}

TEST(IsTest, twoItemsNoCycle) {
  std::shared_ptr<Node> list =
      std::make_shared<Node>(1, std::make_shared<Node>(2));

  const bool got = Is(list);

  EXPECT_FALSE(got);
}

TEST(IsTest, twoItemsCycleToFirst) {
  std::shared_ptr<Node> list =
      std::make_shared<Node>(1, std::make_shared<Node>(2));
  list->next->next = list;

  const bool got = Is(list);

  EXPECT_TRUE(got);
}

TEST(IsTest, twoItemsCycleToSecond) {
  std::shared_ptr<Node> list =
      std::make_shared<Node>(1, std::make_shared<Node>(2));
  list->next->next = list->next;

  const bool got = Is(list);

  EXPECT_TRUE(got);
}

}  // namespace
}  // namespace iq::list::singly::cycle
