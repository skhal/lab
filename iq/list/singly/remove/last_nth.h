// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_LIST_SINGLY_REMOVE_LAST_NTH_H_
#define IQ_LIST_SINGLY_REMOVE_LAST_NTH_H_

#include <cstddef>
#include <memory>

namespace iq::list::singly::remove {

struct Node {
  int value;
  std::shared_ptr<Node> next;
};

std::shared_ptr<Node> RemoveLastNth(std::shared_ptr<Node> list, std::size_t n);

}  // namespace iq::list::singly::remove

#endif  // IQ_LIST_SINGLY_REMOVE_LAST_NTH_H_
