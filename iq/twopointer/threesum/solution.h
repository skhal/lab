// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_TWOPOINTER_THREESUM_SOLUTION_H_
#define IQ_TWOPOINTER_THREESUM_SOLUTION_H_

#include <tuple>
#include <vector>

namespace iq::twopointer::threesum {

using Triplet = std::tuple<int, int, int>;

std::vector<Triplet> Find(const std::vector<int>& nn);

}  // namespace iq::twopointer::threesum

#endif  // IQ_TWOPOINTER_THREESUM_SOLUTION_H_
