# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

# force %c to represent skipped path components with `...`, csh(1)
set ellipsis
# ref: https://en.wikipedia.org/wiki/ANSI_escape_code
# keep-sorted start
set esc_cyan_ = "%{\033[;36m%}"
set esc_green_ = "%{\033[;32m%}"
set esc_reset_ = "%{\033[m%}"
set esc_yellow_ = "%{\033[;33m%}"
# keep-sorted end
set os_ = (`/bin/sh -c '. /etc/os-release; echo $ID $VERSION_ID'`)
set prompt = "${esc_yellow_}${os_}${esc_reset_} ${esc_green_}%N${esc_reset_}@%m:${esc_cyan_}%c03${esc_reset_} %# "
set promptchars = "%#"
unset esc_cyan_ esc_green_ esc_reset_ esc_yellow_ os_
# -- VI(1) PROMPT MODE
bindkey -v
bindkey '^R' i-search-back
