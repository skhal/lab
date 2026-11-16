# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from collections.abc import Callable
from typing import Any, Sequence


class AnsibleModuleProxy:
    """...

    Do not propagate AnsibleModule.params. Pass individual parameter from the
    top-level runner using arguments for explicit data dependency.
    """

    def __init__(
        self,
        *,
        # keep-sorted start
        check_mode: bool,
        run_command: Callable[[Sequence[str]], tuple[int, str, str]],
        # keep-sorted end
    ):
        # keep-sorted start
        self._check_mode = check_mode
        self._run_command = run_command
        # keep-sorted end

    @property
    def check_mode(self):
        return self._check_mode

    def run_command(self, cmd: Sequence[str]) -> tuple[int, str, str]:
        return self._run_command(cmd)
