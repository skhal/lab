# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

# Pick up system install
set base_path_ = "/usr/local/go125"
if ( ! $?TMUX && -d ${base_path_} ) then
  set path = (${path} ${base_path_}/bin)
endif
unset base_path_

# Pick up user binaries from `go install`.
set user_base_path_ = "$HOME"/go
if ( ! $?TMUX && -d ${user_base_path_} ) then
  set path = (${path} ${user_base_path_}/bin)
endif
unset user_base_path_
