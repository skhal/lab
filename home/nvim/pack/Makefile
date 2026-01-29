# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

XDG_CONFIG_HOME ?= ${HOME}/.config
NVIM_CONFIG_PATH = ${XDG_CONFIG_HOME}/nvim
NVIM_PACK_PATH = ${NVIM_CONFIG_PATH}/pack

PLUGINS_FILE = plugins.txt
PLUGINS = ${grep -v '^\#' ${PLUGINS_FILE} | grep '\S':L:sh}

M_target = C,^.*github.com/,,:C,/,/start/,:S,^,${NVIM_PACK_PATH}/,
INSTALL_TARGETS = ${PLUGINS:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

${NVIM_PACK_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for plugin in ${PLUGINS}
target=${plugin:${M_target}}
${target:H}: ${NVIM_PACK_PATH}
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	fi

${target}: ${.TARGET:H} ${PLUGINS_FILE}
	@if ! test -e ${.TARGET}; then \
		echo ".. install ${.TARGET}"; \
		git -C ${.TARGET:H} clone ${plugin}; \
	 fi
.endfor
