// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/10/chain/find.h"

#include <unordered_set>
#include <vector>

namespace iq::chain {
namespace {

std::vector<int> makeChain(const std::unordered_set<int>& nn, int n) {
  std::vector<int> chain = {n};
  for (++n; nn.contains(n); ++n) {
    chain.push_back(n);
  }
  return chain;
}

}  // namespace

std::vector<int> Find(const std::vector<int>& nn) {
  if (nn.empty()) {
    return {};
  }
  std::unordered_set<int> seen{nn.begin(), nn.end()};
  std::vector<int> longest_chain;
  for (const int n : seen) {
    if (const int p = n - 1; seen.contains(p)) {
      continue;
    }
    if (std::vector<int> chain = makeChain(seen, n);
        chain.size() > longest_chain.size()) {
      longest_chain = std::move(chain);
    }
  }
  return longest_chain;
}

}  // namespace iq::chain
