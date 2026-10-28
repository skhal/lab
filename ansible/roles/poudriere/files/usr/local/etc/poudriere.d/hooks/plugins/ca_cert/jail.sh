# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e

[ "$1" == "start" ] || exit 0
[ "${CA_CERT:+set}" == "set" ] || exit 0
[ "${MASTERMNT:+set}" == "set" ] || {
  echo "MASTERMNT is not available" >&2
  exit 1
}

echo "Install CA certificate ${CA_CERT}"

mkdir -p "${MASTERMNT}/${CA_CERT%/*}"
cp "${CA_CERT}" "${MASTERMNT}/${CA_CERT}"

/usr/sbin/certctl -D ${MASTERMNT} rehash
