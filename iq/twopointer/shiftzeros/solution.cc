// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build nobuild

#include <cstddef>
#include <utility>
#include <vector>

#include "iq/twopointer/shiftzeros/solution.h"

namespace iq::twopointer::shiftzeros {
namespace {

constexpr int kZero = 0;

std::size_t indexZero(std::vector<int> nn) {
  std::size_t i = 0;
  while (i < nn.size() && nn[i] != kZero) {
    ++i;
  }
  return i;
}

std::size_t indexNonZero(std::vector<int> nn, std::size_t offset) {
  std::size_t i = offset;
  while (i < nn.size() && nn[i] == kZero) {
    ++i;
  }
  return i;
}

}  // namespace

void Shift(std::vector<int>& nn) {
  std::size_t zero_idx = indexZero(nn);
  if (zero_idx == nn.size()) {
    return;
  }
  for (std::size_t non_zero_idx = indexNonZero(nn, zero_idx + 1);
       non_zero_idx < nn.size();
       non_zero_idx = indexNonZero(nn, non_zero_idx + 1)) {
    std::swap(nn[zero_idx], nn[non_zero_idx]);
    ++zero_idx;
  }
}

}  // namespace iq::twopointer::shiftzeros
