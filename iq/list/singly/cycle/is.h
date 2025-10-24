// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_LIST_SINGLY_CYCLE_IS_H_
#define IQ_LIST_SINGLY_CYCLE_IS_H_

#include <memory>

namespace iq::list::singly::cycle {

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

bool Is(std::shared_ptr<Node> node);

}  // namespace iq::list::singly::cycle

#endif  // IQ_LIST_SINGLY_CYCLE_IS_H_
