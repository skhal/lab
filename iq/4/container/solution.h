// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_CONTAINER_SOLUTION_H_
#define IQ_CONTAINER_SOLUTION_H_

#include <vector>

namespace iq::container {

enum class Volume : int {};

Volume Find(const std::vector<int>& nn);

}  // namespace iq::container

#endif  // IQ_CONTAINER_SOLUTION_H_
