// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/7/twosum/find.h"

#include <cstddef>
#include <optional>
#include <unordered_map>
#include <utility>
#include <vector>

namespace iq::twosum {
namespace {

using Index = std::size_t;

}  // namespace

std::optional<Indices> Find([[maybe_unused]] const std::vector<int>& nn,
                            [[maybe_unused]] int x) {
  std::unordered_map<int, Index> index_by_num;
  for (Index i = 0; i < nn.size(); ++i) {
    const int n = nn[i];
    const int target = x - n;
    if (const auto it = index_by_num.find(target); it != index_by_num.end()) {
      return std::make_pair(it->second, i);
    }
    index_by_num.emplace(n, i);
  }
  return std::nullopt;
}

}  // namespace iq::twosum
