# Copyright (c) 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

VIM_PACK_PATH = ${HOME}/.vim/pack
PLUGIN_LIST = plugins.txt
PLUGINS = ${grep -v '^\#' ${PLUGIN_LIST}:L:sh:M?*}

M_target = C,^.*github.com/,,:S,/,/start/,:S,^,${VIM_PACK_PATH}/,
INSTALL_TARGETS = ${PLUGINS:${M_target}}

install: ${INSTALL_TARGETS} .PHONY

${VIM_PACK_PATH}: .NOTMAIN
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

.for plugin in ${PLUGINS}
target=${plugin:${M_target}}
${target}: ${VIM_PACK_PATH} ${PLUGIN_LIST}
	@if ! test -d ${.TARGET}; then \
		echo ".. install ${plugin} to ${.TARGET}"; \
		mkdir -p ${.TARGET:H} \
			&& cd ${.TARGET:H} \
			&& echo "[+] git clone ${plugin}" \
			&& git clone -q ${plugin} \
			&& echo "[+] vim generate help tags" \
			&& vim -c 'helptags ${.TARGET}/doc' -c 'exit' \
		; \
	 fi
.endfor
