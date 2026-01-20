// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/24/insert/find.h"

#include <cstddef>
#include <vector>

namespace iq::insert {

std::size_t FindInsertIndex(const std::vector<int>& nn, int n) {
  std::size_t left = 0;
  std::size_t right = nn.size();
  while (left < right) {
    const std::size_t pivot = left + (right - left) / 2;
    if (n < nn[pivot]) {
      right = pivot;
    } else if (n > nn[pivot]) {
      left = pivot + 1;
    } else {
      return pivot;
    }
  }
  return left;
}

}  // namespace iq::insert
