# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


class Error:
    def __init__(self, msg: str) -> None:
        self._msg = msg

    def __str__(self) -> str:
        return self._msg
