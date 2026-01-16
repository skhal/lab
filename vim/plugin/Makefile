# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Flavor: FreeBSD

VIM_PLUGIN_PATH = ${HOME}/.vim/plugin
SOURCES = ${find . -maxdepth 1 -type f -name \*.vim:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${VIM_PLUGIN_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${VIM_PLUGIN_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

SOURCE = ${.TARGET:T}

INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}"; \
		cp ${SOURCE} ${.TARGET}

.for target in ${INSTALL_TARGETS}
${target}: ${VIM_PLUGIN_PATH} ${SOURCE} INSTALL_TARGET
.endfor
