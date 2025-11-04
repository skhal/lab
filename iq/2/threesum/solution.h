// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_THREESUM_SOLUTION_H_
#define IQ_THREESUM_SOLUTION_H_

#include <tuple>
#include <vector>

namespace iq::threesum {

using Triplet = std::tuple<int, int, int>;

std::vector<Triplet> Find(const std::vector<int>& nn);

}  // namespace iq::threesum

#endif  // IQ_THREESUM_SOLUTION_H_
