# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

set history = 10000
set savehist = (10000 merge)
set histfile = "${HOME}"/.history

# update histfile if in tmux(1) session and the session had HISTFILE set.
if ( ${?TMUX} ) then
  # Pick up histfile from tmux(1)
  set histfile_ = `tmux show-env HISTFILE >& /dev/stdout`
  if ( ! $? ) then
    set histfile = "${histfile_:s/HISTFILE=//}"
  endif
  unset histfile_
endif

# Isolate Linux VMs history from BSD (default)
if ( "$OSTYPE" == "linux" ) then
  set histfile = "${histfile}.${OSTYPE}"
endif
