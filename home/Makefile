# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Flavor: FreeBSD

.if ${.MAKE.OS} == "Darwin"
PACKAGES = brew

# Following paths use one of the binaries installed by OS packaging system
go npm rust: ${PACKAGES} .NOTMAIN
.endif

INSTALL_PATHS = tcshrc.d nvim ${PACKAGES} ctags.d go npm rust
SOURCES = tcshrc tmux.conf gitconfig
INSTALL_TARGETS = ${SOURCES:S,^,${HOME}/.,}

install: ${INSTALL_TARGETS} ${INSTALL_PATHS} .PHONY

# Target is `${HOME}/.foo`, the source is `foo`.
SOURCE = ${.TARGET:T:S/^.//}

# Instructions to install a target.
INSTALL_TARGET: .NOTMAIN .USE
	@echo ".. install ${SOURCE} to ${.TARGET}" \
		&& cp ${SOURCE} ${.TARGET}

# Generate rules to install targets from sources using the INSTALL_TARGET
# commands.
.for target in ${INSTALL_TARGETS}
${target}: ${SOURCE} INSTALL_TARGET
.endfor

.for target in ${INSTALL_PATHS}
${target}: .PHONY
	@${MAKE} -C ${.CURDIR}/${.TARGET} install
.endfor
