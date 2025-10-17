# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

CTAGS_D_PATH = ${HOME}/.ctags.d
SOURCES = ${find . -maxdepth 1 -type f -name \*.ctags:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${CTAGS_D_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${CTAGS_D_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

SOURCE = ${.TARGET:T}

INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}"; \
		cp ${SOURCE} ${.TARGET}

.for target in ${INSTALL_TARGETS}
${target}: ${CTAGS_D_PATH} ${SOURCE} INSTALL_TARGET
.endfor
