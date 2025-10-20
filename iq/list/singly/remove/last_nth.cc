// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/remove/last_nth.h"

#include <cstddef>
#include <memory>
#include <utility>

namespace iq::list::singly::remove {
namespace {

std::pair<std::shared_ptr<Node>, bool> findPrevToLastNth(
    std::shared_ptr<Node> list, std::size_t n) {
  std::shared_ptr<Node> prev;
  for (std::shared_ptr<Node> node = list; node; node = node->next) {
    if (n > 0) {
      --n;
      continue;
    }
    if (prev == nullptr) {
      prev = list;
      continue;
    }
    prev = prev->next;
  }
  return std::make_pair(prev, n < 1);
}

std::shared_ptr<Node> removeHead(std::shared_ptr<Node> list) {
  std::shared_ptr<Node> next = list->next;
  list->next.reset();
  return next;
}

void removeNextNode(std::shared_ptr<Node> prev) {
  std::shared_ptr<Node> node = prev->next;
  node->next.swap(prev->next);
  node->next.reset();
}

}  // namespace

std::shared_ptr<Node> RemoveLastNth(std::shared_ptr<Node> list, std::size_t n) {
  auto [prev, ok] = findPrevToLastNth(list, n);
  if (!ok) {
    return list;
  }
  if (prev == nullptr) {
    list = removeHead(list);
  } else {
    removeNextNode(prev);
  }
  return list;
}

}  // namespace iq::list::singly::remove
