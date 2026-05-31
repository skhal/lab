# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

# keep-sorted start
GOROOT = ${go env GOROOT:L:sh}
GOVERSION = ${go env GOVERSION:L:sh:C,\.[0-9]+$,,}
# keep-sorted end

M_zipsrc = C,.gz$$,,

install: index.html.gz solsys.js.gz solsys.wasm.gz wasm_exec.js.gz .PHONY

watch: .NOTMAIN .PHONY
	@echo ".. monitor ${.CURDIR} for changes"; \
		make -C ${.CURDIR} \
			&& { \
				fswatch -0 -o ${.CURDIR} \
					| xargs -0 -n 1 -I {} make -C ${.CURDIR}; \
				}

COMPRESS_TARGET: .NOTMAIN .USE
	@echo ".. compress ${.TARGET:${M_zipsrc}}"; \
		gzip --force --keep ${.TARGET:${M_zipsrc}}

index.html.gz: index.html COMPRESS_TARGET

wasm_exec.js.gz: wasm_exec.js COMPRESS_TARGET

wasm_exec.js:
	@if [ -e "${GOROOT}/lib/wasm/${.TARGET}" ]; then \
			echo ".. copy ${.TARGET}" from ${GOROOT}; \
			cp ${GOROOT}/lib/wasm/${.TARGET} ./; \
		else \
			echo ".. fetch ${.TARGET}" from Go repository; \
			fetch -o ./${.TARGET} --quiet https://raw.githubusercontent.com/golang/go/refs/heads/release-branch.${GOVERSION}/lib/wasm/${.TARGET}; \
		fi

solsys.wasm.gz: solsys.wasm COMPRESS_TARGET

solsys.wasm: *.go
	@echo ".. compile ${.TARGET}"; \
		env GOOS=js GOARCH=wasm go build -o ${.TARGET} .

solsys.js.gz: solsys.js COMPRESS_TARGET

solsys.js: solsys.ts
	@echo ".. compile ${.TARGET}"; \
		tsgo ${.ALLSRC}
