// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_INTERSECTION_FIND_H_
#define IQ_INTERSECTION_FIND_H_

#include <memory>

namespace iq::intersection {

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

std::shared_ptr<Node> Find(std::shared_ptr<Node> a, std::shared_ptr<Node> b);

}  // namespace iq::intersection

#endif  // IQ_INTERSECTION_FIND_H_
