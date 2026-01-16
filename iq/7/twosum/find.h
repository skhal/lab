// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_TWOSUM_FIND_H_
#define IQ_TWOSUM_FIND_H_

#include <cstddef>
#include <optional>
#include <utility>
#include <vector>

namespace iq::twosum {

using Indices = std::pair<std::size_t, std::size_t>;

std::optional<Indices> Find(const std::vector<int>& nn, int x);

}  // namespace iq::twosum

#endif  // IQ_TWOSUM_FIND_H_
