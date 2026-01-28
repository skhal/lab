# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

INSTALL_PATHS = lua

install: ${INSTALL_PATHS} .PHONY

.for path in ${INSTALL_PATHS}
${path}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
