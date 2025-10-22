// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_LIST_SINGLY_PALINDROME_IS_H_
#define IQ_LIST_SINGLY_PALINDROME_IS_H_

#include <memory>
namespace iq::list::singly::palindrome {

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

bool Is(const std::shared_ptr<Node>& list);

}  // namespace iq::list::singly::palindrome

#endif  // IQ_LIST_SINGLY_PALINDROME_IS_H_
