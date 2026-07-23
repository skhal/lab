# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

set npm_prefix_bin_ = "$HOME"/npm/node_modules/.bin
if ( ! $?TMUX && -d ${npm_prefix_bin_} ) then
  set path = (${path} ${npm_prefix_bin_})
endif
unset npm_prefix_bin_
