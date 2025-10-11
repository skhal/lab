// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_MAPSET_TWOSUM_FIND_H_
#define IQ_MAPSET_TWOSUM_FIND_H_

#include <cstddef>
#include <optional>
#include <utility>
#include <vector>

namespace iq::mapset::twosum {

using Indices = std::pair<std::size_t, std::size_t>;

std::optional<Indices> Find(const std::vector<int>& nn, int x);

}  // namespace iq::mapset::twosum

#endif  // IQ_MAPSET_TWOSUM_FIND_H_
