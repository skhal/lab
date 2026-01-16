# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Flavor: FreeBSD

CARGO_BIN_PATH = ${HOME}/.cargo/bin
PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST} | tr ' ' ',':L:sh:M*}

M_target = C|,.*$$||:S,^,${CARGO_BIN_PATH}/,
M_unpack = C|,| |g
INSTALL_TARGETS = ${PACKAGES:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

${CARGO_BIN_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for package in ${PACKAGES}
target=${package:${M_target}}
${target}: ${CARGO_BIN_PATH} ${PACKAGE_LIST}
	@if ! test -e ${.TARGET}; then \
		echo ".. install ${.TARGET:T}"; \
		cargo install --locked ${package:${M_unpack}}; \
	 fi
.endfor
