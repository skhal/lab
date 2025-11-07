#!/bin/sh
#
# Copyright 2025 Samvel Khalatyan. All rights reserved.

set -e

readonly GO_VERSION=${GO_VERSION:-1.25.3}
readonly GO_ROOT_PREFIX=/usr/local
readonly LOG_PATH=go_install.log # relative to TMPDIR

CHECKSUM_FILE=$(realpath go_checksum_sha256)
readonly CHECKSUM_FILE

err() {
  echo "$@" >&2
}

check_platform() {
  case $(uname -s) in
    *Linux*) ;;
    *)
      return 1
      ;;
  esac
  case $(uname -v) in
    *FreeBSD*) ;;
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
  return 1
}

make_url() {
  local os=${1:?'make_url: missing os'}
  echo "https://go.dev/dl/go${GO_VERSION}.${os}-amd64.tar.gz"
}

download() {
  local os=${1:?'install: missing os'}
  local url
  url=$(make_url ${os})
  echo ".. download ${url}"
  {
    wget --quiet ${url}
    shasum -c ${CHECKSUM_FILE} --ignore-missing
  } >>${LOG_PATH}
}

extract() {
  local os=${1:?'exract: missing os'}
  local url
  url=$(make_url ${os})
  local archive
  archive=$(basename ${url})
  echo ".. extract ${archive}"
  {
    tar -xzf ${archive}
  } >>${LOG_PATH}
}

patch_from() {
  local os=${1:?'exract: missing os'}
  local url
  url=$(make_url ${os})
  local archive
  archive=$(basename ${url})
  local link_freebsd_path=go/pkg/tool/freebsd_amd64/link
  local link_linux_path=go/pkg/tool/linux_amd64/link
  echo ".. patch ${link_linux_path}"
  {
    tar -xzf ${archive} ${link_freebsd_path}
  } >>${LOG_PATH}
  mv ${link_freebsd_path} ${link_linux_path}
  rm -rf ${link_freebsd_path%/*}
}

make_goroot() {
  local version=${GO_VERSION%.*}              # MAJOR.MINOR
  local suffix=$(echo ${version} | tr -d '.') # {MAJOR}{MINOR}
  echo ${GO_ROOT_PREFIX}/go${suffix}
}

install() {
  local goroot
  goroot=$(make_goroot)
  echo ".. install ${goroot} [need sudo]"
  sudo mv -i ${PWD}/go ${goroot}
}

run() {
  local tmpdir
  tmpdir=$(mktemp -d)
  trap "rm -rf ${tmpdir}" EXIT
  cd ${tmpdir}
  download freebsd
  download linux
  extract linux
  patch_from freebsd
  install
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
