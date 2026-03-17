# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

# force %c to represent skipped path components with `...`, csh(1)
set ellipsis

# ref: https://en.wikipedia.org/wiki/ANSI_escape_code
set esc_green_ = "%{\033[;32m%}"
set esc_reset_ = "%{\033[m%}"
set os_ = (`/bin/sh -c '. /etc/os-release; echo $ID $VERSION_ID'`)
set prompt = "${os_} ${esc_green_}%N${esc_reset_}@%m:%c03 %# "
set promptchars = "%#"
unset esc_green_ esc_reset_
unset os_

# -- VI(1) PROMPT MODE
bindkey -v
bindkey '^R' i-search-back

alias precmd '\\
  set rc_ = $?; \\
  if ( ${rc_} != 0 ) \\
    printf "[\033[;31mrc: ${rc_}\033[0m]\n"; \\
  unset rc_; \\
  '
