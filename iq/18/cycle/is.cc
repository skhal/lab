// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/18/cycle/is.h"

#include <memory>

namespace iq::cycle {

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

}  // namespace iq::cycle
