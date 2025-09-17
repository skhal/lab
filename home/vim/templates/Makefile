# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

VIM_TEMPLATES_PATH = ${HOME}/.vim/templates
SOURCES = ${find . -maxdepth 1 -type f -name skel\*:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:S,^,${VIM_TEMPLATES_PATH}/,}

install: ${INSTALL_TARGETS} .PHONY

${VIM_TEMPLATES_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

SOURCE = ${.TARGET:T}

INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}"; \
		cp ${SOURCE} ${.TARGET}

.for target in ${INSTALL_TARGETS}
${target}: ${VIM_TEMPLATES_PATH} ${SOURCE} INSTALL_TARGET
.endfor
