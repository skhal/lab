# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

GIT_WORKTREE = ${git rev-parse --show-toplevel:L:sh:M*}

PROJECTS = ansible/lab
INSTALL_PROJECTS = ${PROJECTS:S,^,${GIT_WORKTREE}/,:S,$,/.collections,}

ANSIBLE_LOG_PATH = ${HOME}/.ansible/log
ANSIBLE_VAULT_PATH = ${HOME}/.ansible/vault
ANSIBLE_PATHS = ${ANSIBLE_LOG_PATH} ${ANSIBLE_VAULT_PATH}

install: ${ANSIBLE_PATHS} ansible-core ${INSTALL_PROJECTS} .PHONY

CREATE_PATH: .NOTMAIN .USE
	@if ! test -d ${.TARGET}; then \
		echo ".. create ${.TARGET}"; \
		mkdir -p ${.TARGET}; \
	 fi

${ANSIBLE_LOG_PATH}: CREATE_PATH

${ANSIBLE_VAULT_PATH}:: CREATE_PATH

${ANSIBLE_VAULT_PATH}::
	@if test $$(stat -f '%Sp' ${.TARGET}) != "drwx------"; then \
		echo ".. fix permissions for ${.TARGET}"; \
		chmod 0700 ${.TARGET}; \
	 fi
	# unmatched glob stays literal in sh(1): test for file existence to skip it.
	@for f in ${.TARGET}/*; do \
		if ! test -e "$${f}"; then \
			continue; \
		fi; \
		if test $$(stat -f '%Sp' "$$f") != "-r--------"; then \
			echo ".. fix permissions for $$f"; \
			chmod 0400 "$$f"; \
		fi; \
	 done

ansible-core: .PHONY
	@echo ".. install ${.TARGET}"
	@poetry -C ${GIT_WORKTREE} install --no-root

REQUIREMENTS = ${.TARGET:H:S,$,/requirements.yml,}

INSTALL_PROJECT: .NOTMAIN .USE
	@echo ".. install ${REQUIREMENTS} to ${.TARGET}"
	@poetry -C ${.TARGET:H} run \
		ansible-galaxy collection install -r ${REQUIREMENTS} -p ${.TARGET}

.for project in ${INSTALL_PROJECTS}
${project}: ${REQUIREMENTS} INSTALL_PROJECT
.endfor
