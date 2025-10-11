// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build ignore

#include <optional>
#include <utility>
#include <vector>

namespace iq::twopointer::twosum {
namespace {

constexpr int kSizeMin = 2;

}  // namespace

std::optional<std::pair<int, int>> Find(const std::vector<int>& nn, int x) {
  if (nn.size() < kSizeMin) {
    return std::nullopt;
  }
  int leftidx = 0;
  int rightidx = nn.size() - 1;
  while (leftidx < rightidx) {
    const int sum = nn[leftidx] + nn[rightidx];
    if (sum < x) {
      ++leftidx;
    } else if (sum > x) {
      --rightidx;
    } else {
      return std::make_pair(leftidx, rightidx);
    }
  }
  return std::nullopt;
}

}  // namespace iq::twopointer::twosum
