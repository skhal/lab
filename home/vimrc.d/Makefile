# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Flavor: FreeBSD

VIMRC_D_PATH = ${HOME}/.vimrc.d
SOURCES = ${find . -maxdepth 1 -type f -name \*.vim:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${VIMRC_D_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${VIMRC_D_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

SOURCE = ${.TARGET:T}

INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}"; \
		cp ${SOURCE} ${.TARGET}

.for target in ${INSTALL_TARGETS}
${target}: ${VIMRC_D_PATH} ${SOURCE} INSTALL_TARGET
.endfor
