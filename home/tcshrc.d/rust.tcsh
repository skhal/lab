# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

# Pick up user binaries from CARGO_HOME
set prefix_ = "$HOME"/.cargo
if ( ! $?TMUX && -d ${prefix_} ) then
  set path = (${path} ${prefix_}/bin)
endif
unset prefix_
