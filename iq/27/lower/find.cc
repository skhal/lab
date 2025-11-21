// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/27/lower/find.h"

#include <vector>

#include "absl/status/statusor.h"

namespace iq::lower {

absl::StatusOr<int> Find(const std::vector<int>& nn, int k) {
  std::size_t left = 0;
  std::size_t right = nn.size();
  while (left < right) {
    const std::size_t idx = left + (right - left) / 2;
    if (const int n = nn[idx]; k < n) {
      right = idx;
    } else if (k > n) {
      left = idx + 1;
    } else {
      return n;
    }
  }
  if (right == nn.size()) {
    return absl::OutOfRangeError("no such entry");
  }
  return nn[right];
}

}  // namespace iq::lower
