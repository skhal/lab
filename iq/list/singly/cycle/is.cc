// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/cycle/is.h"

#include <memory>

namespace iq::list::singly::cycle {

bool Is([[maybe_unused]] std::shared_ptr<Node> node) {
  std::shared_ptr<Node> slow = node;
  for (std::shared_ptr<Node> fast = node; fast && fast->next;
       fast = fast->next->next) {
    slow = slow->next;
    if (fast.get() == slow.get()) {
      return true;
    }
  }
  return false;
}

}  // namespace iq::list::singly::cycle
