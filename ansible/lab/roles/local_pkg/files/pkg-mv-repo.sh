#!/bin/sh
#
# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# usage: pkg-mv-repo.sh SOURCE TARGET
#
# pkg-mv-repo updates repository annotation in installed packages.
#
# The repository annotation records the repository used to install the packages
# last time and is preferred for the further updates. This script allows to
# switch the repository for installed packages.

usage() {
  cat <<eof >&2
usage: $0 SOURCE TARGET

$0 updates repository annotation in all packages that have SOURCE repository
to the value of TARGET.
eof
}

run() {
  if command -p "${NOOP}"; then
    echo $@
    return
  fi
  eval "$@"
}

NOOP=false
while getopts "n" opt; do
  case "$opt" in
    n)
      NOOP=true
      export NOOP
      ;;
  esac
done
shift $((OPTIND - 1))

if [ $# != 2 ]; then
  usage
  exit 1
fi

SRC_REPO="$1"
TGT_REPO="$2"

pkg annotate --all --show repository \
  | grep "${SRC_REPO}" \
  | sed -E 's/^([^:]+):.*$/\1/' \
  | {
    while read package; do
      run doas pkg annotate --yes --modify $package repository "${TGT_REPO}"
    done
  }
