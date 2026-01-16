// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IQ_GEOMSEQ_FIND_H_
#define IQ_GEOMSEQ_FIND_H_

#include <iosfwd>
#include <vector>

namespace iq::geomseq {

struct Triplet {
  int i, j, k;

  friend std::ostream& operator<<(std::ostream& os, const Triplet& t);
};

std::vector<Triplet> Find(const std::vector<int>& nn, int ratio);

}  // namespace iq::geomseq

#endif  // IQ_GEOMSEQ_FIND_H_
