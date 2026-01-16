// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_STRIPZERO_CLEAN_H_
#define IQ_STRIPZERO_CLEAN_H_

#include <vector>

namespace iq::stripzero {

using Row = std::vector<int>;
using Matrix = std::vector<Row>;

void Clean(Matrix& m);

}  // namespace iq::stripzero

#endif  // IQ_STRIPZERO_CLEAN_H_
