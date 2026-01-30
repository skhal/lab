// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/12/reverse/reverse.h"

#include <memory>

namespace iq::reverse {

std::shared_ptr<Node> Reverse(std::shared_ptr<Node> list) {
  std::shared_ptr<Node> head;
  for (; list; list = list->Next) {
    head = std::make_shared<Node>(list->Value, head);
  }
  return head;
}

}  // namespace iq::reverse
