# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Flavor: FreeBSD

INSTALL_PATHS = lab

VIM_AUTOLOAD_PATH = ${HOME}/.vim/autoload
SOURCES = ${find . -maxdepth 1 -type f -name \*.vim:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${VIM_AUTOLOAD_PATH}/,}

install: ${INSTALL_TARGETS} ${INSTALL_PATHS} .PHONY

${VIM_AUTOLOAD_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

SOURCE = ${.TARGET:T}

INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}"; \
		cp ${SOURCE} ${.TARGET}

.for target in ${INSTALL_TARGETS}
${target}: ${VIM_AUTOLOAD_PATH} ${SOURCE} INSTALL_TARGET
.endfor

.for target in ${INSTALL_PATHS}
${target}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
