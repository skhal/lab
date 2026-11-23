# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

LOCAL_PATH = ${HOME}/.local
LOCAL_BIN_PATH= ${LOCAL_PATH}/bin
LOCAL_SHARE_PATH = ${LOCAL_PATH}/share

VENV_PATH = ${LOCAL_SHARE_PATH}/poetry

install: ${LOCAL_BIN_PATH}/poetry

${LOCAL_BIN_PATH}:
	@echo ".. create ${.TARGET}" && mkdir ${.TARGET}

${VENV_PATH}:
	@echo ".. initialize virtual environment: ${.TARGET:T}" && \
		python3.15 -m venv ${.TARGET}

${LOCAL_BIN_PATH}/poetry: ${LOCAL_BIN_PATH} ${VENV_PATH}
	@echo ".. install poetry" && \
		${VENV_PATH}/bin/pip install poetry && \
		ln -s ${VENV_PATH}/bin/poetry ${.TARGET}
