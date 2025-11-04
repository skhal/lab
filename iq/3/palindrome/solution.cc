// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/3/palindrome/solution.h"

#include <cctype>
#include <string_view>

namespace iq::palindrome {
namespace {

constexpr int kSizeMin = 1;

}  // namespace

bool Is(std::string_view s) {
  if (s.size() < kSizeMin) {
    return true;
  }
  for (int i = 0, j = s.size() - 1; i < j; ++i, --j) {
    while (i < j && !std::isalnum(s[i])) {
      ++i;
    }
    while (i < j && !std::isalnum(s[j])) {
      --j;
    }
    if (i < j && s[i] != s[j]) {
      return false;
    }
  }
  return true;
}

}  // namespace iq::palindrome
