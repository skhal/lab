# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

# -- nvim (linux)
set prefix_ = "/usr/local/nvim-linux-x86_64"
if ( ! $?TMUX && -d ${prefix_} ) then
  set path = (${path} ${prefix_}/bin)
endif
unset prefix_
