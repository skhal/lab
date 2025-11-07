// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_LRU_CACHE_H_
#define IQ_LRU_CACHE_H_

#include <memory>
#include <optional>
#include <vector>

namespace iq::lru {

class Cache {
 public:
  using Key = int;
  using Value = int;

  Cache(std::size_t capacity);

  void Put(Key key, Value value);

  std::optional<Value> Get(Key key);

  std::vector<Key> Keys() const;

 private:
  struct Node {
    Key key;
    Value val;
    std::shared_ptr<Node> next;
    std::shared_ptr<Node> prev;
  };

  bool empty() const;

  std::shared_ptr<Node> find(Key key) const;

  void makeMostRecent(std::shared_ptr<Node> node);

  void removeLeastRecent();

  std::size_t capacity_;
  std::size_t size_ = 0;
  std::shared_ptr<Node> head_;  // most recent
  std::shared_ptr<Node> tail_;  // least recent
};

}  // namespace iq::lru

#endif  // IQ_LRU_CACHE_H_
