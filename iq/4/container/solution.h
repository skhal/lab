// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_CONTAINER_SOLUTION_H_
#define IQ_CONTAINER_SOLUTION_H_

#include <vector>

namespace iq::container {

enum class Volume : int {};

Volume Find(const std::vector<int>& nn);

}  // namespace iq::container

#endif  // IQ_CONTAINER_SOLUTION_H_
