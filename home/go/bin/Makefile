# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

GO_INSTALL_OPTS = -ldflags "-linkmode 'external' -extldflags '-static'"

GO_BIN_PATH = ${HOME}/go/bin
PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

M_target = T:C,@.*$$,,:S,^,${GO_BIN_PATH}/,
INSTALL_TARGETS = ${PACKAGES:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

${GO_BIN_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for package in ${PACKAGES}
target=${package:${M_target}}
${target}: ${GO_BIN_PATH} ${PACKAGE_LIST}
	@if ! test -e ${.TARGET}; then \
		echo ".. install ${package:T}"; \
		go install ${GO_INSTALL_OPTS} ${package}; \
	 fi
.endfor
