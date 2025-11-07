// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_MIDPOINT_FIND_H_
#define IQ_MIDPOINT_FIND_H_

#include <memory>

namespace iq::midpoint {

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

std::shared_ptr<Node> Find(std::shared_ptr<Node> list);

}  // namespace iq::midpoint

#endif  // IQ_MIDPOINT_FIND_H_
