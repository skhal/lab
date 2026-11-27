# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""The exceptions for the patch plugin.

The errors derive from the PatchError.
"""


class PatchError(Exception):
    """Base exception to the poetry patch plugin errors."""


class InvalidSpecError(PatchError):
    """The patch spec is invalid.

    Args:
        msg: an error description.
        spec: a patch spec from TOML configuration.
    """

    def __init__(self, msg: str, spec: dict) -> None:  # noqa: D107
        super().__init__(f"invalid spec: {msg}\n{spec}")
        self.spec = spec
