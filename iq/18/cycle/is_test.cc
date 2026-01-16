// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/18/cycle/is.h"

#include <memory>

#include "gtest/gtest.h"

namespace iq::cycle {
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
}  // namespace iq::cycle
