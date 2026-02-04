// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <iostream>

#include "absl/strings/str_join.h"

int main() { std::cout << absl::StrJoin({"hello", "world\n"}, " "); }
