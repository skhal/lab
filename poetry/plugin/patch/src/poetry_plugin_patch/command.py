# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""The implementation of the `poetry patch` command.

The command reads the `[[tool.poetry-plugin-patch.spec]]` table from the
pyproject.toml and applies patches according to the specs in the configuration
order.

All errors subclass from the PatchError. The command prints errors to the
standard error and reports non-zero exit code.
"""

from functools import cached_property
from pathlib import Path
from typing import override

from poetry.console.commands.env_command import EnvCommand

from poetry_plugin_patch._error import format_exc
from poetry_plugin_patch.error import (
    AlreadyPatchedError,
    PackageNotFoundError,
    PatchError,
)
from poetry_plugin_patch.patch import patch
from poetry_plugin_patch.spec import PatchSpec


class PatchCommand(EnvCommand):
    """It implements `poetry patch` command."""

    name = "patch"
    description = "Patch the installed packages."
    help = """\
The <info>patch</info> command patches installed packages in the project
environment according to the patch declarations in the pyproject.toml.

[[tool.poetry-plugin-patch.spec]]
package = "example"
patches = ["path/to/file.patch"]
"""

    @cached_property
    def pyproject_base(self) -> Path:
        """The root path of the Python project."""

        return self.poetry.pyproject_path.parent

    @cached_property
    def specs(self) -> list[dict]:
        """A list of patch specifications in the TOML config."""

        return (
            self.poetry.pyproject.data.get("tool", {})
            .get("poetry-plugin-patch", {})
            .get("spec", [])
        )

    @override
    def handle(self) -> int:
        """Apply patches from TOML specification.

        Returns:
            0: on success
            1: on failure
        """

        if not self.specs:
            return 0

        try:
            self._handle()
        except PatchError as err:
            self.line_error(format_exc(err))
            return 1

        return 0

    def _handle(self) -> None:
        specs = tuple(
            PatchSpec.parse(s, self.pyproject_base) for s in self.specs
        )

        for spec in specs:
            dist = self.env.site_packages.find_distribution(spec.package)
            if dist is None:
                raise PackageNotFoundError(spec.package)

            try:
                patch(dist, spec)
            except AlreadyPatchedError as err:
                self.line(str(err))
            else:
                self.line(f"patched {spec.package}")
