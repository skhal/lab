# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# .cshrc - csh resource script, read at beginning of execution by each shell
#
# see also csh(1), environ(7).
# more examples available at /usr/share/examples/csh/
#

alias h		history 25
alias j		jobs -l
alias la	ls -aF
alias lf	ls -FA
alias ll	ls -lAF

# read(2) of directories may not be desirable by default, as this will provoke
# EISDIR errors from each directory encountered.
# alias grep	grep -d skip

# A righteous umask
umask 22

set path = (/sbin /bin /usr/sbin /usr/bin /usr/local/sbin /usr/local/bin $HOME/bin)

setenv	EDITOR	vi
setenv	PAGER	less

if ($?prompt) then
	# An interactive shell -- set some stuff up
	# force %c to represent skipped path components with `...`, csh(1)
	set ellipsis
	# ANSI colors:
  # - 3x foreground
  # - 4x background
	#   0 black
	#   1 red
	#   2 green
	#   3 yellow
	#   4 blue
	#   5 magenta
	#   6 cyan
	# ref: https://en.wikipedia.org/wiki/ANSI_escape_code
	set esc_color_ = "%{\033[;31m%}"
	set esc_reset_ = "%{\033[m%}"
	set os_ = (`/bin/sh -c '. /etc/os-release; echo $ID $VERSION_ID'`)
	set prompt = "${os_} ${esc_color_}%N${esc_reset_}@%m:%c03 %# "
	set promptchars = "%#"
	unset esc_color_ esc_reset_ os_

	set filec
	set history = 1000
	set savehist = (1000 merge)
	set autolist = ambiguous
	# Use history to aid expansion
	set autoexpand
	set autorehash
	set mail = (/var/mail/$USER)
	if ( $?tcsh ) then
		bindkey "^W" backward-delete-word
		bindkey -k up history-search-backward
		bindkey -k down history-search-forward
	endif

endif
