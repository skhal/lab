// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_ROTATED_FIND_H_
#define IQ_ROTATED_FIND_H_

#include <vector>

#include "absl/status/statusor.h"

namespace iq::rotated {

absl::StatusOr<int> Find(const std::vector<int>& nn, int k);

}  // namespace iq::rotated

#endif  // IQ_ROTATED_FIND_H_
