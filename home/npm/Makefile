# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

NPM_PATH=${HOME}/npm
BIN_PATH=${NPM_PATH}/node_modules/.bin
DOAS_CMD=/usr/local/bin/doas
.if ${.MAKE.OS} == Darwin
DOAS_CMD=/usr/bin/sudo
.endif

NPM_BIN=${NPM_PATH}/bin

PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

M_target = C|@.*$$||:S,^,${BIN_PATH}/,
M_no_target = C|^.*@|@|
INSTALL_TARGETS = ${PACKAGES:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

${NPM_PATH}: .NOTMAIN
	@echo ".. create ${.TARGET}"; \
		mkdir ${.TARGET}

.for package in ${PACKAGES}
target=${package:${M_target}}
${target}: ${PACKAGE_LIST} ${NPM_PATH}
	@echo ".. install ${.TARGET:T}"; \
		npm install --no-fund --silent --prefix ${NPM_PATH} ${package:${M_no_target}};
.endfor
