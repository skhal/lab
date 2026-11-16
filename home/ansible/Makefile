# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

GIT_WORKTREE = ${git rev-parse --show-toplevel:L:sh:M*}

PROJECTS = ansible/lab
INSTALL_PROJECTS = ${PROJECTS:S,^,${GIT_WORKTREE}/,:S,$,/.collections,}

install:  ansible-core ${INSTALL_PROJECTS} .PHONY

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
