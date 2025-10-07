// Copyright 2025 Samvel Khalatyan. All rights reserved.

#include <algorithm>

#include "iq/twopointer/container/solution.h"

namespace iq::twopointer::container {
namespace {

constexpr Volume kVolumeEmpty = Volume{0};

constexpr int kSizeMin = 2;

Volume CalculateVolume(const std::vector<int>& nn, int lidx, int ridx) {
  const int height = std::min(nn[lidx], nn[ridx]);
  const int width = ridx - lidx;
  return Volume{height * width};
}

}  // namespace

Volume Find(const std::vector<int>& nn) {
  if (nn.size() < kSizeMin) {
    return kVolumeEmpty;
  }
  int lidx = 0;
  int ridx = nn.size() - 1;
  Volume vol_max = kVolumeEmpty;
  while (lidx < ridx) {
    if (const Volume vol = CalculateVolume(nn, lidx, ridx); vol > vol_max) {
      vol_max = vol;
    }
    if (const int ln = nn[lidx], rn = nn[ridx]; ln < rn) {
      lidx += 1;
    } else if (ln > rn) {
      ridx -= 1;
    } else {
      lidx += 1;
      ridx -= 1;
    }
  }
  return vol_max;
}

}  // namespace iq::twopointer::container
