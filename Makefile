# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

INSTALL_TARGETS = home nvim vim

install: ${INSTALL_TARGETS} .PHONY

# Install home after vim.
#
# `make -C home` installs vimrc and Vim plugins. It runs Vim to generate help
# tags. The installed vimrc uses a colorscheme and plugins installed by
# `make -C vim`.
home: vim .PHONY

.for target in ${INSTALL_TARGETS}
${target}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
