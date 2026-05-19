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
	# 0 black
	# 1 red
	# 2 green
	# 3 yellow
	# 4 blue
	# 5 magenta
	# 6 cyan
	# ref: https://en.wikipedia.org/wiki/ANSI_escape_code
	set esc_red_ = "%{\033[;31m%}"
	set esc_green_ = "%{\033[;32m%}"
        set esc_yellow_ = "%{\033[;1;33m%}"
        set esc_blue_ = "%{\033[;34m%}"
        set esc_magenta_ = "%{\033[;35m%}"
	set esc_cyan_ = "%{\033[;36m%}"
        set esc_reset_ = "%{\033[m%}"
	set os_ = (`/bin/sh -c '. /etc/os-release; echo $ID $VERSION_ID'`)
	set prompt = "${os_} ${esc_yellow_}%N${esc_reset_}@%m:%c03 %# "
	set promptchars = "%#"
	unset esc_red_ esc_green_ esc_yellow_ esc_blue_ esc_magenta_ esc_cyan_ esc_reset_
	unset os_

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
