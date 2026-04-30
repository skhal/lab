# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Flavor: FreeBSD

GOVERSION = ${go env GOVERSION:L:sh:C,\.[0-9]+$,,}

install: hello.wasm .PHONY

hello.wasm: wasm_exec.js *.go
	@echo ".. compile ${.TARGET}"; \
		env GOOS=js GOARCH=wasm go build -o ${.TARGET} .

wasm_exec.js:
	@echo ".. fetch ${.TARGET}"; \
		fetch --quiet https://raw.githubusercontent.com/golang/go/refs/heads/release-branch.${GOVERSION}/lib/wasm/${.TARGET}
