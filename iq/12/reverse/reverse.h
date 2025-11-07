// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_REVERSE_REVERSE_H_
#define IQ_REVERSE_REVERSE_H_

#include <memory>

namespace iq::reverse {

struct Node {
  int Value;
  std::shared_ptr<Node> Next;
};

std::shared_ptr<Node> Reverse(std::shared_ptr<Node> list);

}  // namespace iq::reverse

#endif  // IQ_REVERSE_REVERSE_H_
