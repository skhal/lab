// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <iostream>
#include <string>
#include <vector>

#include "absl/strings/str_join.h"

int main() {
  std::vector<std::string> v = {"foo", "bar", "baz"};
  std::string s = absl::StrJoin(v, "-");

  std::cout << "Joined string: " << s << "\n";

  return 0;
}
