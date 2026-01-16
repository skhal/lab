// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_REMOVE_LAST_NTH_H_
#define IQ_REMOVE_LAST_NTH_H_

#include <cstddef>
#include <memory>

namespace iq::remove {

struct Node {
  int value;
  std::shared_ptr<Node> next;
};

std::shared_ptr<Node> RemoveLastNth(std::shared_ptr<Node> list, std::size_t n);

}  // namespace iq::remove

#endif  // IQ_REMOVE_LAST_NTH_H_
