# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

BIN_PATH=/usr/local/bin
DOAS_CMD=/usr/local/bin/doas
.if ${.MAKE.OS} == Darwin
DOAS_CMD=/usr/bin/sudo
.endif

PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

M_target = C|@.*$$||:S,^,${BIN_PATH}/,
M_no_target = C|^.*@|@|
INSTALL_TARGETS = ${PACKAGES:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

.for package in ${PACKAGES}
target=${package:${M_target}}
${target}: ${PACKAGE_LIST}
	@if which ${.TARGET:T} >/dev/null; then \
		echo ".. skip ${.TARGET:T}"; \
	 else \
		echo ".. install ${.TARGET:T}" to ${.TARGET}; \
		${DOAS_CMD} npm install --global ${package:${M_no_target}}; \
	 fi
.endfor
