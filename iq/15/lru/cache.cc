// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/15/lru/cache.h"

#include <memory>
#include <optional>
#include <vector>

#include "absl/base/macros.h"

namespace iq::lru {

Cache::Cache(std::size_t capacity) : capacity_(capacity) {
  ABSL_ASSERT(capacity > 0);
}

void Cache::Put(Key key, Value value) {
  if (auto node = find(key); node) {
    node->val = value;
    makeMostRecent(std::move(node));
    return;
  }
  if (!empty() && size_ == capacity_) {
    removeLeastRecent();
  }
  auto node = std::make_shared<Node>(key, value, head_);
  if (!empty()) {
    head_->prev = node;
  } else {
    tail_ = node;
  }
  head_ = std::move(node);
  ++size_;
}

std::optional<Cache::Value> Cache::Get(Key key) {
  const std::shared_ptr<Node> node = find(key);
  if (!node) {
    return std::nullopt;
  }
  makeMostRecent(node);
  return node->val;
}

std::vector<Cache::Key> Cache::Keys() const {
  std::vector<Cache::Key> keys;
  for (std::shared_ptr<Node> node = head_; node; node = node->next) {
    keys.emplace_back(node->key);
  }
  return keys;
}

bool Cache::empty() const { return !head_; }

std::shared_ptr<Cache::Node> Cache::find(Key key) const {
  std::shared_ptr<Node> node = head_;
  for (; node; node = node->next) {
    if (node->key == key) {
      break;
    }
  }
  return node;
}

void Cache::makeMostRecent(std::shared_ptr<Node> node) {
  if (!node->prev) {
    // first node
    return;
  }
  if (!node->next) {
    // last node
    tail_ = node->prev;
    tail_->next.reset();
    node->prev.reset();
  } else {
    auto prev = node->prev;
    auto next = node->next;
    prev->next = next;
    next->prev = prev;
    node->next.reset();
    node->prev.reset();
  }
  node->next = head_;
  head_->prev = node;
  head_ = std::move(node);
}

void Cache::removeLeastRecent() {
  auto prev = tail_->prev;
  if (!prev) {
    tail_.reset();
    head_.reset();
  } else {
    prev->next.reset();
    tail_->prev.reset();
    tail_ = std::move(prev);
  }
  --size_;
}

}  // namespace iq::lru
