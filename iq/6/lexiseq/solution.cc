// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/6/lexiseq/solution.h"

#include <algorithm>
#include <iterator>
#include <string>
#include <string_view>
#include <tuple>
#include <utility>

namespace iq::lexiseq {
namespace {

constexpr int kSizeMin = 2;

struct Pivot {
  char ch;
  int idx;

  friend bool operator==(const Pivot& a, const Pivot& b) {
    return std::tie(a.ch, a.idx) == std::tie(b.ch, b.idx);
  }
};

constexpr Pivot kPivotInvalid = {.ch = -1, .idx = -1};

// findPivot scans the string s backwards from the end until it finds the first
// non-ascending character, called pivot. It returns the found character and
// its index, or kPivotInvalid if the search fails.
Pivot findPivot(std::string_view s) {
  for (int i = s.size() - 1; i > 0; --i) {
    if (const int p = i - 1; s[p] < s[i]) {
      return {.ch = s[p], .idx = p};
    }
  }
  return kPivotInvalid;
}

// findPivotNext scans the string s bakwards from the end until it finds the
// first character that is larger than the pivot. It returns the position of
// the character or zero if the search fails.
int findPivotNext(std::string_view s, char pivot) {
  for (int i = s.size() - 1; i >= 0; --i) {
    if (s[i] > pivot) {
      return i;
    }
  }
  return 0;
}

}  // namespace

std::string Next(std::string_view s) {
  std::string buf{s};
  if (s.size() < kSizeMin) {
    return buf;
  }
  const Pivot pivot = findPivot(buf);
  if (pivot == kPivotInvalid) {
    std::reverse(buf.begin(), buf.end());
    return buf;
  }
  const int pivot_next_idx =
      pivot.idx + 1 + findPivotNext(s.substr(pivot.idx + 1), pivot.ch);
  std::swap(buf[pivot.idx], buf[pivot_next_idx]);
  std::reverse(
      [&buf, &pivot]() {
        auto it = buf.begin();
        std::advance(it, pivot.idx + 1);
        return it;
      }(),
      buf.end());
  return buf;
}

}  // namespace iq::lexiseq
