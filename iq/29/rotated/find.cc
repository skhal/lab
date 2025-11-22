// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/29/rotated/find.h"

#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"

namespace iq::rotated {

absl::StatusOr<int> Find(const std::vector<int>& nn, int k) {
  std::size_t low = 0;
  std::size_t hi = nn.size();
  while (low < hi) {
    const std::size_t idx = low + (hi - low) / 2;
    if (int n = nn[idx]; k < n) {
      // May need to move right due to rotation
      if (k >= nn[low]) {
        hi = idx;
      } else {
        low = idx + 1;
      }
    } else if (k > n) {
      // May need to move left due to rotation
      if (k <= nn[hi - 1]) {
        low = idx + 1;
      } else {
        hi = idx;
      }
    } else {
      return idx;
    }
  }
  return absl::OutOfRangeError("out of range");
}

}  // namespace iq::rotated
