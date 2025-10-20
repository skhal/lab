// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/reverse/reverse.h"

#include <memory>

namespace iq::list::singly::reverse {

std::shared_ptr<Node> Reverse(std::shared_ptr<Node> list) {
  std::shared_ptr<Node> head;
  for (; list; list = list->Next) {
    head = std::make_shared<Node>(list->Value, head);
  }
  return head;
}

}  // namespace iq::list::singly::reverse
