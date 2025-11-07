// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_FLATTEN_DO_H_
#define IQ_FLATTEN_DO_H_

#include <memory>

namespace iq::flatten {

struct TreeNode {
  int val;
  std::shared_ptr<TreeNode> next;
  std::shared_ptr<TreeNode> child;
};

struct Node {
  int val;
  std::shared_ptr<Node> next;
};

std::shared_ptr<Node> Do(std::shared_ptr<TreeNode> tree);

}  // namespace iq::flatten

#endif  // IQ_FLATTEN_DO_H_
