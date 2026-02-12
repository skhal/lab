# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

TCSHRCD_PATH = ${HOME}/.tcshrc.d
SOURCES = ${find . -maxdepth 1 -type f -name \*.tcsh:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${TCSHRCD_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${TCSHRCD_PATH}:
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for target in ${INSTALL_TARGETS}
${target}: ${TCSHRCD_PATH} ${.TARGET:T}
	@echo ".. install ${.TARGET:T} to ${.TARGET}"; \
		cp ${.TARGET:T} ${.TARGET}
.endfor
