// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/16/palindrome/is.h"

#include <cstddef>
#include <memory>
#include <utility>

namespace iq::palindrome {
namespace {

std::pair<std::shared_ptr<Node>, std::shared_ptr<Node>> findMidpoint(
    std::shared_ptr<Node> list) {
  std::shared_ptr<Node> prev;
  std::shared_ptr<Node> next;
  for (std::size_t size = 1; list; list = list->next, size += 1) {
    if (next == nullptr) {
      next = list;
      prev = std::make_shared<Node>(list->val, prev);
      continue;
    }
    if (size % 2 == 0) {
      next = next->next;
    } else {
      prev = std::make_shared<Node>(next->val, prev);
    }
  }
  return std::make_pair(prev, next);
}

}  // namespace

bool Is(const std::shared_ptr<Node>& list) {
  if (list == nullptr) {
    return true;
  }
  if (list->next == nullptr) {
    return true;
  }
  auto [prev, next] = findMidpoint(list);
  while (prev != nullptr && next != nullptr) {
    if (prev->val != next->val) {
      return false;
    }
    prev = prev->next;
    next = next->next;
  }
  return prev == nullptr && next == nullptr;
}

}  // namespace iq::palindrome
