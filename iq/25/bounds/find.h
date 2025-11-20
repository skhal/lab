// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_BOUNDS_FIND_H_
#define IQ_BOUNDS_FIND_H_

#include <utility>
#include <vector>

namespace iq::bounds {

using Bounds = std::pair<int, int>;

inline constexpr int kIllegalIndex = -1;
inline constexpr Bounds kIllegalBounds = {kIllegalIndex, kIllegalIndex};

Bounds Find(const std::vector<int>& nn, int n);

}  // namespace iq::bounds

#endif  // IQ_BOUNDS_FIND_H_
