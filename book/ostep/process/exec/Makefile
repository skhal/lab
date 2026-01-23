# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

TARGET=/tmp/ostep-process-${.CURDIR:T}

${TARGET}: main.c
	@echo ".. build ${.TARGET}" \
		&& cc -o ${.TARGET} ${.ALLSRC}

run: ${TARGET} .PHONY
	@echo ".. run ${.ALLSRC}" \
		&& ${.ALLSRC}
