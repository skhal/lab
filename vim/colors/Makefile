# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

VIM_COLORS_PATH = ${HOME}/.vim/colors
SOURCES = ${find . -maxdepth 1 -type f -name \*.vim:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${VIM_COLORS_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${VIM_COLORS_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

SOURCE = ${.TARGET:T}

INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}"; \
		cp ${SOURCE} ${.TARGET}

.for target in ${INSTALL_TARGETS}
${target}: ${VIM_COLORS_PATH} ${SOURCE} INSTALL_TARGET
.endfor
