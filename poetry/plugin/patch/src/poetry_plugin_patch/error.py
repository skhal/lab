# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""The exceptions for the patch plugin.

The errors derive from the PatchError.
"""

from pathlib import Path


class PatchError(Exception):
    """Base exception to the poetry patch plugin errors."""


class AlreadyPatchedError(PatchError):
    """The package is already patched.

    Args:
        package: the package name.
    """

    def __init__(self, package: str) -> None:  # noqa: D107
        super().__init__(f"package {package} is already patched")
        self.package = package


class CopyError(PatchError):
    """There was error copying package files.

    Args:
        src: the source path.
        dst: the destination path.
    """

    def __init__(self, src: Path, dst: Path) -> None:  # noqa: D107
        super().__init__(f"failed to copy {src} to {dst}")
        self.src = src
        self.dst = dst


class EmptyDistributionRecordsError(PatchError):
    """Distribution has no records.

    Args:
        package: the package name
    """

    def __init__(self, package: str) -> None:  # noqa: D107
        super().__init__(f"distribution {package} has empty records")
        self.package = package


class InvalidPatchError(PatchError):
    """The patch is invalid.

    Args:
        patch: the patch path.
        output: combined standard error and output of the patch(1) run.
    """

    def __init__(self, patch: Path, output: str) -> None:  # noqa: D107
        super().__init__(f"invalid patch {patch}\n{output}")
        self.patch = patch
        self.output = output


class InvalidSpecError(PatchError):
    """The patch spec is invalid.

    Args:
        msg: an error description.
        spec: a patch spec from TOML configuration.
    """

    def __init__(self, msg: str, spec: dict) -> None:  # noqa: D107
        super().__init__(f"invalid spec: {msg}\n{spec}")
        self.spec = spec


class PackageNotFoundError(PatchError):
    """The package is not found.

    Args:
        package: the package name.
    """

    def __init__(self, package: str) -> None:  # noqa: D107
        super().__init__(f"package not found {package}")
        self.package = package


class PackagePatchError(PatchError):
    """Failed to apply patches to the package.

    Args:
        package: the package name
    """

    def __init__(self, package: str) -> None:  # noqa: D107
        super().__init__(f"failed to patch {package}")
        self.package = package


class PatchNotFoundError(PatchError):
    """patch(1) is not available."""

    def __init__(self) -> None:  # noqa: D107
        super().__init__("patch(1) is not available")


class RemoveFolderError(PatchError):
    """Error removing a folder.

    Args:
        path: the folder path.
    """

    def __init__(self, path: Path) -> None:  # noqa: D107
        super().__init__(f"failed to remove the folder {path}")
        self.path = path


class UnlinkError(PatchError):
    """Error removing a file.

    Args:
        file: the file path.
    """

    def __init__(self, file: Path) -> None:  # noqa: D107
        super().__init__(f"failed to unlink the file {file}")
        self.file = file
