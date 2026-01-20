// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_ANAGRAM_FIND_ALL_H_
#define IQ_ANAGRAM_FIND_ALL_H_

#include <string>
#include <string_view>
#include <vector>
namespace iq::anagram {

std::vector<std::string> FindAll(std::string_view s, std::string_view t);

}  // namespace iq::anagram

#endif  // IQ_ANAGRAM_FIND_ALL_H_
