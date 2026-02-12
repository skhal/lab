# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

if ( ${?TMUX} ) then
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
