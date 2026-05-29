# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

M_zipsrc = C,.gz$$,,

install: index.html.gz solsys.css.gz solsys.js.gz .PHONY

watch: .NOTMAIN .PHONY
	@echo ".. monitor ${.CURDIR} for changes"; \
		fswatch -o ${.CURDIR} | xargs -n 1 -I {} make -C ${.CURDIR}

COMPRESS_TARGET: .NOTMAIN .USE
	@echo ".. compress ${.TARGET:${M_zipsrc}}"; \
		gzip --force --keep ${.TARGET:${M_zipsrc}}

index.html.gz: index.html COMPRESS_TARGET

solsys.css.gz: solsys.css COMPRESS_TARGET

solsys.js.gz: solsys.js COMPRESS_TARGET

solsys.js: solsys.ts
	@echo ".. compile ${.TARGET}"; \
		tsgo ${.ALLSRC}
