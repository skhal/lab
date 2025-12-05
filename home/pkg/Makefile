# Copyright (c) 2025 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

INSTALL_TARGETS = ${PACKAGES}

install: ${INSTALL_TARGETS} .PHONY

.for package in ${INSTALL_TARGETS}
${package}: ${PACKAGE_LIST}
	@if ! pkg info -e ${.TARGET}; then \
		echo ".. install ${.TARGET}"; \
		doas pkg install -y ${.TARGET} \
		; \
	 fi
.endfor
