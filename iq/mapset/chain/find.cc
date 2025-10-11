// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build ignore

#include <unordered_set>
#include <vector>

#include "iq/mapset/chain/find.h"

namespace iq::mapset::chain {
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

}  // namespace iq::mapset::chain
