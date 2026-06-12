# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

# keep-sorted start
M_jszip  = C,.ts$$,.js.gz,
M_zipsrc = C,.gz$$,,
# keep-sorted end

SOURCES = ${find . -maxdepth 1 -type f -name \*.ts:L:sh:S,./,,}
INSTALL_TARGETS = ${SOURCES:${M_jszip}}

install: ${INSTALL_TARGETS} .PHONY

.SUFFIXES: .js.gz .js

.js.js.gz:
	@echo ".. compress ${.IMPSRC}"; \
		gzip --force --keep ${.IMPSRC}
