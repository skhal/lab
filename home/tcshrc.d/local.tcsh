# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

if ( -d "$HOME"/.local/bin ) then
  set path = ($path "$HOME"/.local/bin)
  # -f to keep the first occurrence, -l to keep the last one.
  set -f path = ($path)
endif
