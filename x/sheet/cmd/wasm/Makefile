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

COMPRESS_TARGET: .NOTMAIN .USE
	@echo ".. compress ${.TARGET:${M_zipsrc}}"; \
		gzip --force --keep ${.TARGET:${M_zipsrc}}

install: index.html.gz sheet.wasm.gz wasm_exec.js.gz .PHONY

index.html.gz: index.html COMPRESS_TARGET

sheet.wasm.gz: sheet.wasm COMPRESS_TARGET

sheet.wasm: *.go
	@echo ".. compile ${.TARGET}"; \
		env GOOS=js GOARCH=wasm go build -o ${.TARGET} .

wasm_exec.js.gz: wasm_exec.js COMPRESS_TARGET

wasm_exec.js:
	@if [ -e "${GOROOT}/lib/wasm/${.TARGET}" ]; then \
			echo ".. copy ${.TARGET}" from ${GOROOT}; \
			cp ${GOROOT}/lib/wasm/${.TARGET} ./; \
		else \
			echo ".. fetch ${.TARGET}" from Go repository; \
			fetch --quiet https://raw.githubusercontent.com/golang/go/refs/heads/release-branch.${GOVERSION}/lib/wasm/${.TARGET}; \
		fi
