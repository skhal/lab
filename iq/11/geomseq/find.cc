// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/11/geomseq/find.h"

#include <cstddef>
#include <initializer_list>
#include <ostream>
#include <unordered_map>
#include <vector>

#include "absl/strings/str_format.h"

namespace iq::geomseq {
namespace {

using Index = std::size_t;
using Indices = std::vector<Index>;

class PrevNumbers {
 public:
  const Indices* Get(int n) const {
    auto it = index_by_num_.find(n);
    if (it == index_by_num_.end()) {
      return nullptr;
    }
    return &(it->second);
  }

  void Add(int n, Index i) {
    auto it = index_by_num_.find(n);
    if (it == index_by_num_.end()) {
      index_by_num_.emplace(n, std::initializer_list<Index>{i});
      return;
    }
    it->second.push_back(i);
  }

 private:
  std::unordered_map<int, Indices> index_by_num_;
};

class NextNumbers {
 public:
  NextNumbers(const std::vector<int>& nn) {
    for (Index i = nn.size(); i-- > 0;) {
      auto it = index_by_num_.find(nn[i]);
      if (it == index_by_num_.end()) {
        index_by_num_.emplace(nn[i], std::initializer_list<Index>{i});
        return;
      }
      it->second.push_back(i);
    }
  }

  const Indices* Get(int n) const {
    auto it = index_by_num_.find(n);
    if (it == index_by_num_.end()) {
      return nullptr;
    }
    return &(it->second);
  }

  void Pop(int n) {
    auto it = index_by_num_.find(n);
    if (it == index_by_num_.end()) {
      return;
    }
    it->second.pop_back();
  }

 private:
  std::unordered_map<int, Indices> index_by_num_;
};

class Triplets {
 public:
  void Add(const Indices& prev_ii, Index i, const Indices& next_ii) {
    for (Index pi : prev_ii) {
      for (Index ni : next_ii) {
        triplets_.emplace_back(pi, i, ni);
      }
    }
  }

  operator std::vector<Triplet>() { return triplets_; }

 private:
  std::vector<Triplet> triplets_;
};

}  // namespace

std::ostream& operator<<(std::ostream& os, const Triplet& t) {
  return os << absl::StrFormat("[%d %d %d]", t.i, t.j, t.k);
}

std::vector<Triplet> Find(const std::vector<int>& nn, int ratio) {
  PrevNumbers nums_prev;
  NextNumbers nums_next(nn);
  Triplets triplets;
  for (Index i = 0; i < nn.size(); ++i) {
    const int n = nn[i];
    nums_next.Pop(n);
    const Indices* prev_ii = nums_prev.Get(n / ratio);
    const Indices* next_ii = nums_next.Get(n * ratio);
    if (prev_ii != nullptr && next_ii != nullptr) {
      triplets.Add(*prev_ii, i, *next_ii);
    }
    nums_prev.Add(n, i);
  }
  return triplets;
}

}  // namespace iq::geomseq
