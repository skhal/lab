// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/28/upper/find.h"

#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"

namespace iq::upper {

absl::StatusOr<int> Find([[maybe_unused]] const std::vector<int>& nn,
                         [[maybe_unused]] int k) {
  std::size_t low = 0;
  std::size_t hi = nn.size();
  while (low < hi) {
    const std::size_t idx = low + (hi - low) / 2;
    if (const int n = nn[idx]; k < n) {
      hi = idx;
    } else if (k > n) {
      low = idx + 1;
    } else {
      hi = idx + 1;
      break;
    }
  }
  if (hi == nn.size()) {
    return absl::OutOfRangeError("no such entry");
  }
  return nn[hi];
}

}  // namespace iq::upper
