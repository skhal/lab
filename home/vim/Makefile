# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

INSTALL_TARGETS = pack templates

install: ${INSTALL_TARGETS} .PHONY

.for target in ${INSTALL_TARGETS}
${target}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
