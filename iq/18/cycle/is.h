// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_CYCLE_IS_H_
#define IQ_CYCLE_IS_H_

#include <memory>

namespace iq::cycle {

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

bool Is(std::shared_ptr<Node> node);

}  // namespace iq::cycle

#endif  // IQ_CYCLE_IS_H_
