// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/2/threesum/solution.h"

#include <algorithm>
#include <cstddef>
#include <utility>
#include <vector>

namespace iq::threesum {
namespace {

constexpr int kSizeMin = 3;
constexpr int kSizePositiveMin = 2;

enum class Sign : int {
  Negative = 1,
  Positive = 2,
};

Sign sign(int n) {
  if (n >= 0) {
    return Sign::Positive;
  }
  return Sign::Negative;
}

bool hasOppositeSignEnds(const std::vector<int>& nn) {
  const Sign sign_left = sign(*nn.begin());
  const Sign sign_right = sign(*nn.rbegin());
  return sign_left != sign_right;
}

std::size_t indexPositive(const std::vector<int>& nn) {
  std::size_t i = 0;
  while (i < nn.size() && nn[i] < 0) {
    ++i;
  }
  return i;
}

using TwoSum = std::pair<int, int>;

std::vector<TwoSum> findTwoSumAll(const std::vector<int>& nn,
                                  std::size_t offset, int n) {
  std::vector<TwoSum> two_sums;
  std::size_t i = offset;
  std::size_t j = nn.size() - 1;
  while (i < j) {
    const int sum = nn[i] + nn[j];
    if (sum < n) {
      ++i;
    } else if (sum > n) {
      --j;
    } else {
      two_sums.emplace_back(nn[i], nn[j]);
      do {
        ++i;
      } while (i < j && nn[i] == nn[i - 1]);
      do {
        --j;
      } while (i < j && nn[j] == nn[j + 1]);
    }
  }
  return two_sums;
}

std::size_t indexNegative(const std::vector<int>& nn, std::size_t offset,
                          std::size_t size) {
  while (offset < size && nn[offset - 1] == nn[offset]) {
    ++offset;
  }
  return offset;
}

}  // namespace

std::vector<Triplet> Find(const std::vector<int>& nn) {
  if (nn.size() < kSizeMin) {
    return {};
  }
  std::vector<int> tmp{nn};
  std::sort(tmp.begin(), tmp.end());
  if (!hasOppositeSignEnds(tmp)) {
    return {};
  }
  const std::size_t posidx = indexPositive(tmp);
  if (tmp.size() - posidx < kSizePositiveMin) {
    return {};
  }
  std::vector<Triplet> triplets;
  for (std::size_t i = 0; i < posidx; i = indexNegative(tmp, i + 1, posidx)) {
    std::vector<TwoSum> two_sums = findTwoSumAll(tmp, posidx, -tmp[i]);
    for (auto [a, b] : std::move(two_sums)) {
      triplets.emplace_back(tmp[i], a, b);
    }
  }
  return triplets;
}

}  // namespace iq::threesum
