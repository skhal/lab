# Copyright (c) 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

VIM_PACK_PATH = ${HOME}/.vim/pack
PACKAGE_LIST = packages.txt
PACKAGES = ${grep -v '^\#' ${PACKAGE_LIST}:L:sh:M?*}

M_target = C,^.*github.com/,,:S,/,/start/,:S,^,${VIM_PACK_PATH}/,
INSTALL_TARGETS = ${PACKAGES:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

${VIM_PACK_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for package in ${PACKAGES}
target=${package:${M_target}}
${target}: ${VIM_PACK_PATH} ${PACKAGE_LIST}
	@if ! test -d ${.TARGET}; then \
		echo ".. install ${package} to ${.TARGET}"; \
		mkdir -p ${.TARGET:H} \
			&& cd ${.TARGET:H} \
			&& echo "[+] git clone ${package}" \
			&& git clone -q ${package} \
			&& echo "[+] vim generate help tags" \
			&& vim -c 'helptags ${.TARGET}/doc' -c 'exit' \
		; \
	 fi
.endfor
