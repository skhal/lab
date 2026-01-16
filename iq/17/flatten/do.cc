// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// clang-format off-next-line
//go:build ignore

#include "iq/17/flatten/do.h"

#include <functional>
#include <memory>
#include <queue>
#include <utility>

namespace iq::flatten {
namespace {

std::pair<std::shared_ptr<Node>, std::shared_ptr<Node>> toList(
    std::shared_ptr<TreeNode> tree,
    std::function<void(std::shared_ptr<TreeNode>)> visit_child) {
  std::shared_ptr<Node> head;
  std::shared_ptr<Node> tail;
  for (; tree; tree = tree->next) {
    if (tree->child) {
      visit_child(tree->child);
    }
    auto node = std::make_shared<Node>(tree->val);
    if (!head) {
      head = node;
    } else {
      tail->next = node;
    }
    tail = std::move(node);
  }
  return std::make_pair(std::move(head), std::move(tail));
}

}  // namespace

std::shared_ptr<Node> Do([[maybe_unused]] std::shared_ptr<TreeNode> tree) {
  std::shared_ptr<Node> head;
  std::shared_ptr<Node> tail;
  std::queue<std::shared_ptr<TreeNode>> queue;
  auto enque = [&queue](std::shared_ptr<TreeNode> tree_node) {
    queue.emplace(std::move(tree_node));
  };
  for (enque(std::move(tree)); !queue.empty(); queue.pop()) {
    std::shared_ptr<TreeNode> tree_node = queue.front();
    auto [h, t] = toList(tree_node, enque);
    if (!head) {
      head = std::move(h);
    } else {
      tail->next = std::move(h);
    }
    tail = std::move(t);
  }
  return head;
}

}  // namespace iq::flatten
