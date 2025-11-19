// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// clang-format off-next-line
//go:build ignore

#include "iq/21/anagram/find_all.h"

#include <cstddef>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

namespace iq::anagram {
namespace {

class Footprint {
 public:
  explicit Footprint(std::string_view s) {
    for (auto ch : s) {
      footprint_[ch] += 1;
    }
  }

  bool operator==(const Footprint& other) const {
    if (this == &other) {
      return true;
    }
    return footprint_ == other.footprint_;
  }

 private:
  std::unordered_map<char, std::size_t> footprint_;
};

}  // namespace

std::vector<std::string> FindAll(std::string_view s, std::string_view t) {
  if (s.empty() || s.size() < t.size()) {
    return {};
  }
  std::vector<std::string> anagrams;
  const Footprint want_footprint(t);
  for (std::size_t i = 0, tlen = t.size(); i + tlen <= s.size(); ++i) {
    const std::string_view substr = s.substr(i, tlen);
    Footprint substr_footprint(substr);
    if (want_footprint != substr_footprint) {
      continue;
    }
    anagrams.emplace_back(substr);
  }
  return anagrams;
}

}  // namespace iq::anagram
