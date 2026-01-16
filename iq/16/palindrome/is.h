// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
