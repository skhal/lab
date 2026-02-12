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
  # Pick up ssh-agent socket from tmux(1)
  alias tmux-update-environment '\\
    if ( ${?TMUX} ) \\
      setenv SSH_AUTH_SOCK `tmux show-env SSH_AUTH_SOCK | cut -d'=' -f 2` \\
    '
  set postcmd_ = `alias postcmd \
    | sed -E 's/tmux-update-environment//' \
    | sed -E 's/;[[:space:]]*;/;/' \
    | sed -E 's/;$//'`
  if ( ${?postcmd_} && ${#postcmd_} ) then
    alias postcmd "${postcmd_}; tmux-update-environment"
  else
    alias postcmd tmux-update-environment
  endif
  unset postcmd_
endif

# Isolate Linux VMs history from BSD (default)
if ( "$OSTYPE" == "linux" ) then
  set histfile = "${histfile}.${OSTYPE}"
endif
