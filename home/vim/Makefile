# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Flavor: FreeBSD

INSTALL_TARGETS = pack templates

install: ${INSTALL_TARGETS} .PHONY

.for target in ${INSTALL_TARGETS}
${target}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
