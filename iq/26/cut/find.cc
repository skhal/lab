// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/26/cut/find.h"

#include <vector>

namespace iq::cut {
namespace {

int max(const std::vector<int>& nn) {
  int max_n = nn.front();
  for (const int n : nn) {
    if (n > max_n) {
      max_n = n;
    }
  }
  return max_n;
}

int sumAbove(const std::vector<int>& nn, int cut) {
  int sum = 0;
  for (const int n : nn) {
    if (n > cut) {
      sum += (n - cut);
    }
  }
  return sum;
}

}  // namespace

int Find(const std::vector<int>& nn, int k) {
  int low = 0;
  int hi = max(nn);
  while (low < hi) {
    const int mid = low + (hi - low) / 2;
    const int sum = sumAbove(nn, mid + 1);  // +1 to ceil()
    if (sum < k) {
      hi = mid;
    } else {
      low = mid + 1;
    }
  }
  return low;
}

}  // namespace iq::cut
