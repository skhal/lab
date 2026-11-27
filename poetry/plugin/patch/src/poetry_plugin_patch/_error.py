# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


def format_exc(err: BaseException | None) -> str:
    lines = []
    while err is not None:
        lines.append(str(err))
        err = err.__cause__
    return "\n".join(lines)
