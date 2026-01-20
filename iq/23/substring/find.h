// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_SUBSTRING_FIND_H_
#define IQ_SUBSTRING_FIND_H_

#include <cstddef>
#include <string_view>

namespace iq::substring {

std::string_view Find(std::string_view s, std::size_t n);

}  // namespace iq::substring

#endif  // IQ_SUBSTRING_FIND_H_
