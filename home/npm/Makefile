# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

BIN_PATH=/usr/local/bin

PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

M_target = C|@.*$$||:S,^,${BIN_PATH}/,
M_no_target = C|^.*@|@|
INSTALL_TARGETS = ${PACKAGES:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

.for package in ${PACKAGES}
target=${package:${M_target}}
${target}: ${PACKAGE_LIST}
	echo ".. install ${.TARGET:T}" to ${.TARGET}; \
		doas npm install --global ${package:${M_no_target}}
.endfor
