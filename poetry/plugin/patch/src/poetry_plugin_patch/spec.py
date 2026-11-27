# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""It is a parsed specification from the plugin's TOML configuration.

A PatchSpec is a single, validated entry:

    [[tool.poetry-plugin-patch.spec]]
    package = "name"
    patches = ["one.patch", "path/two.patch"]
    strip = 2  # optional

Use PatchSpec.parse() to validate the spec. It raises InvalidSpecError when the
validation fails.
"""

from dataclasses import dataclass, fields
from pathlib import Path

from poetry_plugin_patch.error import InvalidSpecError


@dataclass(frozen=True)
class PatchSpec:
    """It represents a single item in the plugin's TOML configuration."""

    package: str
    patches: tuple[Path, ...]
    strip: int = 1

    @classmethod
    def parse(cls, spec: dict, base: Path) -> "PatchSpec":
        """Validate spec params and create a new PatchSpec.

        The patch files must be relative to the pyproject in the spec but
        PatchSpec keeps the files as absolute paths with respect to base.

        Args:
            spec: a patch specification in TOML format.
            base: root folder of the patch files.

        Returns:
            A PatchSpec object.

        Raises:
            InvalidSpecError: when spec parameters are invalid.
        """

        names = {f.name for f in fields(cls)}
        if unknown := sorted(set(spec) - names):
            raise InvalidSpecError(
                f"unsupported fields {', '.join(unknown)}", spec
            )

        package = spec.get("package")
        if not package or not isinstance(package, str):
            raise InvalidSpecError("missing package name", spec)

        patches = spec.get("patches")
        if (
            patches is None
            or not isinstance(patches, list)
            or not all(isinstance(p, str) for p in patches)
        ):
            raise InvalidSpecError("patches must be a list of strings", spec)

        rel_patches = tuple(Path(p) for p in patches)
        for p in rel_patches:
            if p.is_absolute():
                raise InvalidSpecError(f"patch must be relative - {p}", spec)

        abs_patches = tuple(base / p for p in rel_patches)
        for abs, rel in zip(abs_patches, rel_patches, strict=True):
            if not abs.is_file():
                raise InvalidSpecError(f"patch not found - {rel}", spec)

        strip = spec.get("strip", 1)
        if not isinstance(strip, int) or isinstance(strip, bool):
            raise InvalidSpecError("strip must be an integer", spec)
        if strip < 0:
            raise InvalidSpecError("strip can not be negative", spec)

        return cls(package, abs_patches, strip)
