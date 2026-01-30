// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/5/shiftzeros/solution.h"

#include <cstddef>
#include <utility>
#include <vector>

namespace iq::shiftzeros {
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

}  // namespace iq::shiftzeros
