// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/25/bounds/find.h"

#include <functional>
#include <vector>

namespace iq::bounds {
namespace {

using FindFunc = std::function<int(const std::vector<int>&, int, std::size_t,
                                   std::size_t, std::size_t)>;

int find(const std::vector<int>& nn, int n, std::size_t left, std::size_t right,
         FindFunc fn) {
  while (left < right) {
    const std::size_t idx = left + (right - left) / 2;
    if (n < nn[idx]) {
      right = idx;
    } else if (n > nn[idx]) {
      left = idx + 1;
    } else {
      return fn(nn, n, left, idx, right);
    }
  }
  return kIllegalIndex;
}

int findLeft(const std::vector<int>& nn, int n, std::size_t left,
             std::size_t pivot, [[maybe_unused]] std::size_t right) {
  if (const int idx = find(nn, n, left, pivot, findLeft);
      idx != kIllegalIndex) {
    return idx;
  }
  return pivot;
}

int findRight(const std::vector<int>& nn, int n,
              [[maybe_unused]] std::size_t left, std::size_t pivot,
              std::size_t right) {
  if (const int idx = find(nn, n, pivot + 1, right, findRight);
      idx != kIllegalIndex) {
    return idx;
  }
  return pivot;
}

}  // namespace

Bounds Find(const std::vector<int>& nn, int n) {
  const int left = find(nn, n, 0, nn.size(), findLeft);
  if (left == kIllegalIndex) {
    return kIllegalBounds;
  }
  const int right = find(nn, n, left, nn.size(), findRight);
  if (right == kIllegalIndex) {
    return kIllegalBounds;
  }
  return {left, right};
}

}  // namespace iq::bounds
