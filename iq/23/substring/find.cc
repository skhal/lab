// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "iq/23/substring/find.h"

#include <algorithm>
#include <cstddef>
#include <string_view>
#include <unordered_map>

namespace iq::substring {

std::string_view Find(std::string_view s, std::size_t n) {
  std::string_view substr;
  std::unordered_map<char, std::size_t> freq;
  std::size_t freq_max = 0;  // maximum frequency of the last valid window
  std::size_t start = 0;
  std::size_t end = 0;
  while (end < s.size()) {
    const char ch = s[end];
    end += 1;
    freq[ch] += 1;
    freq_max = std::max(freq_max, freq[ch]);
    if (const std::size_t replacements = (end - start) - freq_max;
        replacements > n) {
      const char ch = s[start];
      freq[ch] -= 1;
      start += 1;
    }
    if (const std::size_t size = end - start; size > substr.size()) {
      substr = s.substr(start, size);
    }
  }
  return substr;
}

}  // namespace iq::substring
