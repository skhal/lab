// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line: skip file in `go build`
//go:build ignore

#include <cstddef>

#include "iq/mapset/stripzero/clean.h"

namespace iq::mapset::stripzero {
namespace {

constexpr int kZero = 0;

struct CleanOptions {
  bool first_row_has_zero = false;
  bool first_col_has_zero = false;
};

CleanOptions flagMatrix(Matrix& m) {
  CleanOptions opts;
  for (std::size_t r = 0; r < m.size(); ++r) {
    Row& row = m[r];
    for (std::size_t c = 0; c < row.size(); ++c) {
      if (row[c] != kZero) {
        continue;
      }
      if (r == 0) {
        opts.first_row_has_zero = true;
      }
      if (c == 0) {
        opts.first_col_has_zero = true;
      }
      m[r][0] = kZero;
      m[0][c] = kZero;
    }
  }
  return opts;
}

void cleanRow(Row& row) {
  for (std::size_t c = 0; c < row.size(); ++c) {
    row[c] = kZero;
  }
}

void cleanRowsButFirst(Matrix& m) {
  for (std::size_t r = 1; r < m.size(); ++r) {
    Row& row = m[r];
    if (row[0] != kZero) {
      continue;
    }
    cleanRow(row);
  }
}

void cleanCol(Matrix& m, std::size_t c) {
  for (Row& row : m) {
    row[c] = kZero;
  }
}

void cleanColsButFirst(Matrix& m) {
  const Row& first_row = m[0];
  for (std::size_t c = 1; c < first_row.size(); ++c) {
    if (first_row[c] != kZero) {
      continue;
    }
    cleanCol(m, c);
  }
}

void cleanMatrix(Matrix& m, const CleanOptions& opts) {
  cleanRowsButFirst(m);
  cleanColsButFirst(m);
  if (opts.first_row_has_zero) {
    cleanRow(m[0]);
  }
  if (opts.first_col_has_zero) {
    cleanCol(m, 0);
  }
}

}  // namespace

void Clean(Matrix& m) {
  if (m.empty()) {
    return;
  }
  const CleanOptions opts = flagMatrix(m);
  cleanMatrix(m, opts);
}

}  // namespace iq::mapset::stripzero
