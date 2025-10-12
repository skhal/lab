// Copyright 2025 Samvel Khalatyan. All rights reserved.

#ifndef IQ_MAPSET_GEOMSEQ_FIND_H_
#define IQ_MAPSET_GEOMSEQ_FIND_H_

#include <iosfwd>
#include <vector>

namespace iq::mapset::geomseq {

struct Triplet {
  int i, j, k;

  friend std::ostream& operator<<(std::ostream& os, const Triplet& t);
};

std::vector<Triplet> Find(const std::vector<int>& nn, int ratio);

}  // namespace iq::mapset::geomseq

#endif  // IQ_MAPSET_GEOMSEQ_FIND_H_
