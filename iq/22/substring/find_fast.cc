// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NOLINTNEXTLINE
//go:build ignore

#include "iq/22/substring/find_fast.h"

#include <cstddef>
#include <unordered_map>

namespace iq::substring {
namespace {

constexpr std::size_t kMinSize = 2;
constexpr char kCharInvalid = -1;

class Window {
 public:
  explicit Window(std::string_view s) : str_(s) {}

  bool Scan() {
    if (end_ == str_.size()) {
      return false;
    }
    shrink();
    grow();
    return true;
  }

  std::size_t Size() const { return end_ - start_; }

  std::string_view String() const { return str_.substr(start_, end_ - start_); }

 private:
  void grow() {
    for (std::size_t i = end_; i < str_.size(); ++i) {
      const char ch = str_[i];
      const auto it = positions_.find(ch);
      if (it == positions_.end()) {
        positions_[ch] = i;
        ++end_;
      } else if (it->second < start_) {
        positions_[ch] = i;
        ++end_;
      } else {
        stop_ = ch;
        break;
      }
    }
  }

  void shrink() {
    if (stop_ == kCharInvalid) {
      return;
    }
    start_ = positions_[stop_] + 1;
    stop_ = kCharInvalid;
  }

 private:
  std::string_view str_;
  std::size_t start_ = 0;
  std::size_t end_ = 0;
  char stop_ = kCharInvalid;
  std::unordered_map<char, std::size_t> positions_;
};

}  // namespace

std::string_view FindFast(std::string_view s) {
  if (s.size() < kMinSize) {
    return s;
  }
  std::string_view substr;
  Window w(s);
  while (w.Scan()) {
    if (w.Size() <= substr.size()) {
      continue;
    }
    substr = w.String();
  }
  return substr;
}

}  // namespace iq::substring
