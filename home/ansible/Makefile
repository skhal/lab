# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Flavor: FreeBSD

# Poetry installs ansible-core in a virtual environment. Collections are not on
# PyPI, so use ansible-galaxy from the environment to install them into the
# default location, ~/.ansible/collections.
ROOT = ${git rev-parse --show-toplevel:L:sh:M*}
REQUIREMENTS = ${ROOT}/ansible/requirements.yml

install: .PHONY
	@echo ".. install ansible-core with poetry"
	@poetry -C ${ROOT} install --no-root
	@echo ".. install ansible collections from ${REQUIREMENTS:S,^${ROOT}/,,}"
	@poetry -C ${ROOT} run ansible-galaxy collection install -r ${REQUIREMENTS}
