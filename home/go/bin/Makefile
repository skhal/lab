# Copyright 2025 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

GO_INSTALL_OPTS = -ldflags "-s -w -linkmode external -extldflags -static"

GO_BIN_PATH = ${HOME}/go/bin
PACKAGE_LIST = packages.txt
# Package line in packages.txt:
#		package [flags]
#
# make(1) treats white-space for field separator, i.e. space, tab, eol, etc. Use
# comma to separate fields in the line to keep package and flags together:
#   package,flag,[flag...]
#
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST} | tr ' ' ',':L:sh:M*}

M_package = C|,.*$$||
M_flags = C|^[^,]*,||
M_unpack = C|,| |g
# Target is the last component of the package, `name`. A package may have an optional
# version suffix:
#   host/path/name@version
#
M_target = ${M_package}:T:C|@.*$$||:S|^|${GO_BIN_PATH}/|
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
		echo ".. install ${.TARGET:T}"; \
		go install ${GO_INSTALL_OPTS} ${package:${M_flags}:${M_unpack}} ${package:${M_package}}; \
	 fi
.endfor
