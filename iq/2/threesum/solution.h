// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_THREESUM_SOLUTION_H_
#define IQ_THREESUM_SOLUTION_H_

#include <tuple>
#include <vector>

namespace iq::threesum {

using Triplet = std::tuple<int, int, int>;

std::vector<Triplet> Find(const std::vector<int>& nn);

}  // namespace iq::threesum

#endif  // IQ_THREESUM_SOLUTION_H_
