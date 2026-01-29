# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

XDG_CONFIG_HOME ?= ${HOME}/.config
NVIM_CONFIG_PATH = ${XDG_CONFIG_HOME}/nvim
TARGET_PATH = ${NVIM_CONFIG_PATH}/lua/config

SOURCES = ${find . -maxdepth 1 -type f -name \*.lua:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${TARGET_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${TARGET_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for target in ${INSTALL_TARGETS}
${target}: ${TARGET_PATH} ${.TARGET:T}
	@echo ".. install ${.TARGET:T} to ${.TARGET}"; \
		cp ${.TARGET:T} ${.TARGET}
.endfor
