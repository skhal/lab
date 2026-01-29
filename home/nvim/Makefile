# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

XDG_CONFIG_HOME ?= ${HOME}/.config
NVIM_CONFIG_PATH = ${XDG_CONFIG_HOME}/nvim

INSTALL_PATHS = after lua pack

SOURCES = init.lua
INSTALL_TARGETS = ${SOURCES:S,^,${NVIM_CONFIG_PATH}/,}

install: ${INSTALL_PATHS} ${INSTALL_TARGETS} .PHONY

${NVIM_CONFIG_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for target in ${INSTALL_TARGETS}
${target}: ${NVIM_CONFIG_PATH} ${.TARGET:T}
	@echo ".. install ${.TARGET:T} to ${.TARGET}"; \
		cp ${.TARGET:T} ${.TARGET}
.endfor

.for target in ${INSTALL_PATHS}
${target}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
