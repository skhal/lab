// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/19/midpoint/find.h"

#include <memory>

namespace iq::midpoint {

std::shared_ptr<Node> Find(std::shared_ptr<Node> list) {
  std::shared_ptr<Node> slow = list;
  for (std::shared_ptr<Node> fast = list; fast && fast->next;
       fast = fast->next->next) {
    slow = slow->next;
  }
  return slow;
}

}  // namespace iq::midpoint
