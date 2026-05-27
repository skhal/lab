# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

M_zipsrc = C,.gz$$,,

fswatch: .NOTMAIN .PHONY
	@echo ".. monitor ${.CURDIR} for changes"; \
		fswatch -o ${.CURDIR} | xargs -n 1 -I {} make -C ${.CURDIR}

COMPRESS_TARGET: .NOTMAIN .USE
	@echo ".. compress ${.TARGET:${M_zipsrc}}"; \
		gzip --force --keep ${.TARGET:${M_zipsrc}}

install: canvas.js.gz index.html.gz .PHONY

canvas.js.gz: canvas.js COMPRESS_TARGET

canvas.js: canvas.ts
	@echo ".. compile ${.TARGET}"; \
		tsgo ${.ALLSRC}

index.html.gz: index.html COMPRESS_TARGET
