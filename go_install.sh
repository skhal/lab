#!/bin/sh
#
# Copyright 2025 Samvel Khalatyan. All rights reserved.

set -e

readonly GO_VERSION=${GO_VERSION:-1.25.3}
readonly GO_ROOT_PREFIX=/usr/local

CHECKSUM_FILE=$(realpath go_checksum_sha256)
readonly CHECKSUM_FILE

err() {
  echo "$@" >&2
}

check_platform() {
  case $(uname -s) in
  *Linux*)
    ;;
  *)
    return 1
    ;;
  esac
  case $(uname -v) in
  *FreeBSD*)
    ;;
  *)
    return 1
    ;;
  esac
}

check_privileges() {
  if [ $(id -u) = "0" ]; then
    return
  fi
  case $(id -nG) in
  *sudoer*)
    return
    ;;
  esac
  return  1
}

download() {
  local os=${1:?'install: missing os'}
  local url="https://go.dev/dl/go${GO_VERSION}.${os}-amd64.tar.gz"
  local archive
  archive=$(basename ${url})
  wget ${url}
  shasum \
    -c ${CHECKSUM_FILE} \
    --ignore-missing
  mkdir ${os}
  tar \
    -C ${os} \
    -xzf ${archive}
}

patch_linker() {
  cp \
    ./freebsd/go/pkg/tool/freebsd_amd64/link \
    ./linux/go/pkg/tool/linux_amd64/link
}

install() {
  local version=${GO_VERSION%.*}
  local goroot=${GO_ROOT_PREFIX}/go$(echo ${version} | tr -d '.')
  sudo mv -i ${PWD}/linux/go ${goroot}
  echo "installed Go ${GO_VERSION} to ${goroot}"
}

run() {
  local tmpdir
  tmpdir=$(mktemp -d)
  trap "rm -rf ${tmpdir}" EXIT
  cd ${tmpdir}
  download freebsd
  download linux
  patch_linker
  install ./linux/go
}

main() {
  if ! check_platform; then
    err "unsupported platform - want Linux compatibility under FreeBSD"
    return 1
  fi
  if ! check_privileges; then
    err "the effective user does not have sudo privileges"
    return 1
  fi
  run
}

main
