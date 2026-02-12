# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

set history = 10000
set savehist = (10000 merge)
# Isolate Linux VMs history from BSD (default)
if ( "$OSTYPE" == "linux" ) then
  set histfile = "$HOME"/.history."$OSTYPE"
endif

# update histfile for tmux(1) sessions if one is set
if ( ${?TMUX} ) then
  # Pick up histfile from tmux(1)
  set histfile_ = `tmux show-env HISTFILE >& /dev/stdout`
  if ( ! $? ) then
    set histfile_ = "${histfile_:s/HISTFILE=//}"
    if ( "$OSTYPE" == "linux" ) then
      set histfile_ = "${histfile_}.$OSTYPE"
    endif
    set histfile = "${histfile_}"
  endif
  unset histfile_
  # Pick up ssh-agent socket from tmux(1)
  alias dev-update-environment '\\
    if ( ${?TMUX} ) \\
      setenv SSH_AUTH_SOCK `tmux show-env SSH_AUTH_SOCK | cut -d'=' -f 2` \\
    '
  set postcmd_ = `alias postcmd \
    | sed -E 's/dev-update-environment//' \
    | sed -E 's/;[[:space:]]*;/;/' \
    | sed -E 's/;$//'`
  if ( ${?postcmd_} && ${#postcmd_} ) then
    alias postcmd "${postcmd_}; dev-update-environment"
  else
    alias postcmd dev-update-environment
  endif
  unset postcmd_
endif
