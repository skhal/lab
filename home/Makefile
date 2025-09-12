# Copyright (c) 2025 Samvel Khalatyan. All rights reserved.
#
# Flavor: FreeBSD

INSTALL_TARGETS = ${HOME}/.tcshrc ${HOME}/.vimrc

install: ${INSTALL_TARGETS} .PHONY

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
