// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_STRIPZERO_CLEAN_H_
#define IQ_STRIPZERO_CLEAN_H_

#include <vector>

namespace iq::stripzero {

using Row = std::vector<int>;
using Matrix = std::vector<Row>;

void Clean(Matrix& m);

}  // namespace iq::stripzero

#endif  // IQ_STRIPZERO_CLEAN_H_
