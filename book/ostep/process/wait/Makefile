# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

${.CURDIR:T}: main.c
	@echo ".. build ${.TARGET}" \
		&& cc -o ${.TARGET} ${.ALLSRC}

run: ${.CURDIR:T} .PHONY
	./${.CURDIR:T}
