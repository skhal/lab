// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_LOWER_FIND_H_
#define IQ_LOWER_FIND_H_

#include <vector>

#include "absl/status/statusor.h"

namespace iq::lower {

absl::StatusOr<int> Find(const std::vector<int>& nn, int k);

}  // namespace iq::lower

#endif  // IQ_LOWER_FIND_H_
