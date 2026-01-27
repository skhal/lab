# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

XDG_CONFIG_HOME ?= ${HOME}/.config
NVIM_CONFIG_PATH = ${XDG_CONFIG_HOME}/nvim
TARGET_PATH = ${NVIM_CONFIG_PATH}/lua/plugins

install: ${TARGET_PATH} .PHONY

${TARGET_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -vp ${.TARGET}; \
	 fi
