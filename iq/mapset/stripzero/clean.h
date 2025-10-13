// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_MAPSET_STRIPZERO_CLEAN_H_
#define IQ_MAPSET_STRIPZERO_CLEAN_H_

#include <vector>

namespace iq::mapset::stripzero {

using Row = std::vector<int>;
using Matrix = std::vector<Row>;

void Clean(Matrix& m);

}  // namespace iq::mapset::stripzero

#endif  // IQ_MAPSET_STRIPZERO_CLEAN_H_
