# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

SOURCES = init.lua
XDG_CONFIG_HOME ?= ${HOME}/.config
NVIM_CONFIG_PATH = ${XDG_CONFIG_HOME}/nvim
INSTALL_TARGETS = ${SOURCES:S,^,${NVIM_CONFIG_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${NVIM_CONFIG_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -vp ${.TARGET}; \
	 fi

.for target in ${INSTALL_TARGETS}
${target}: ${NVIM_CONFIG_PATH} .PHONY
	@echo ".. install ${.TARGET:T} to ${.TARGET}"; \
		cp -v ${.TARGET:T} ${.TARGET}
.endfor
