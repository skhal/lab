// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/14/intersection/find.h"

#include <memory>
#include <tuple>

namespace iq::intersection {
namespace {

struct Chain {
  std::shared_ptr<Node> head;
  std::shared_ptr<Node> tail;

  class ForwardIterator {
   public:
    ForwardIterator(Chain* chain, bool is_begin)
        : chain_(chain),
          is_head_(is_begin),
          node_(is_begin ? chain->head : nullptr) {}

    ForwardIterator& operator++() {
      if (node_) {
        node_ = node_->next;
        if (!node_ && is_head_) {
          is_head_ = false;
          node_ = chain_->tail;
        }
      }
      return *this;
    }

    const std::shared_ptr<Node>& operator*() const { return node_; }

    bool operator==(const ForwardIterator& that) const {
      if (this == &that) {
        return true;
      }
      return std::tie(chain_, is_head_, node_) ==
             std::tie(that.chain_, that.is_head_, that.node_);
    }

   private:
    const Chain* chain_;
    bool is_head_;
    std::shared_ptr<Node> node_;
  };

  ForwardIterator begin() { return ForwardIterator(this, true); }

  ForwardIterator end() { return ForwardIterator(this, false); }
};

std::shared_ptr<Node> find(Chain a, Chain b) {
  auto ait = a.begin();
  auto aend = a.end();
  auto bit = b.begin();
  auto bend = b.end();
  while (ait != aend && bit != bend) {
    if (*ait == *bit) {
      return *ait;
    }
    ++ait;
    ++bit;
  }
  return nullptr;
}

}  // namespace

std::shared_ptr<Node> Find(std::shared_ptr<Node> a, std::shared_ptr<Node> b) {
  return find(Chain(a, b), Chain(b, a));
}

}  // namespace iq::intersection
