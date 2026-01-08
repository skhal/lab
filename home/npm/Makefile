# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

BIN_PATH=/usr/local/bin

PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

INSTALL_TARGETS = ${PACKAGES:S,^,${BIN_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

.for target in ${INSTALL_TARGETS}
${target}: ${PACKAGE_LIST}
	@echo ".. install ${.TARGET:T}"; \
	 doas npm install --global ${.TARGET:T}
.endfor
