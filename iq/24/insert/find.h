// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_INSERT_FIND_H_
#define IQ_INSERT_FIND_H_

#include <cstddef>
#include <vector>

namespace iq::insert {

std::size_t FindInsertIndex(const std::vector<int>& nn, int n);

}  // namespace iq::insert

#endif  // IQ_INSERT_FIND_H_
