// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_PALINDROME_IS_H_
#define IQ_PALINDROME_IS_H_

#include <memory>

namespace iq::palindrome {

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

bool Is(const std::shared_ptr<Node>& list);

}  // namespace iq::palindrome

#endif  // IQ_PALINDROME_IS_H_
