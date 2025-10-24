// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/midpoint/find.h"

#include <memory>

namespace iq::list::singly::midpoint {

std::shared_ptr<Node> Find(std::shared_ptr<Node> list) {
  std::shared_ptr<Node> slow = list;
  for (std::shared_ptr<Node> fast = list; fast && fast->next;
       fast = fast->next->next) {
    slow = slow->next;
  }
  return slow;
}

}  // namespace iq::list::singly::midpoint
