# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""Apply patches to the distribution.

The patch module applies patches to the distribution in a temporary folder and
copies the patched results back. It removes the files if one of the patches
removes it and cleans up corresponding .pyc files.

All errors subclass from PatchError.
"""

import os
import shutil
import subprocess
import tempfile
from collections.abc import Iterable
from contextlib import suppress
from enum import IntEnum, StrEnum, auto
from importlib.metadata import Distribution, PackagePath
from pathlib import Path

from poetry_plugin_patch.error import (
    AlreadyPatchedError,
    CopyError,
    EmptyDistributionRecordsError,
    InvalidPatchError,
    PackagePatchError,
    PatchNotFoundError,
    RemoveFolderError,
    UnlinkError,
)
from poetry_plugin_patch.package import Package
from poetry_plugin_patch.spec import PatchSpec


def patch(dist: Distribution, spec: PatchSpec) -> None:
    """Apply patch files from the spec to the distribution.

    patch() attempts to patch the distribution in a temporary folder and copies
    the results back to the distribution if it succeeds. It removes files from
    the distribution if the patch does so.

    It leaves the distribution in a mixed state in the following scenarios:

    - a copy of the patched files to the distribution fails.
    - a clean up of the removed files from the distribution fails.
    - the distribution RECORD manifest refers the original files: hashes, sizes,
      etc.

    Therefore it becomes important to re-patch the distribution after every
    poetry upgrade or install.

    Args:
        dist: distribution to patch.
        spec: patch specification.

    Raises:
        AlreadyPatchedError: the distribution is already patched.
        CopyError: it fails to copy the distribution to or patched distribution
            from the temporary folder.
        PackagePatchError: one of the patch files is invalid.
        PatchNotFoundError: patch(1) is not available.
        RemoveFolderError: failed to remove a folder from the distribution.
        UnlinkError: failed to remove a file from the distribution.
    """

    with tempfile.TemporaryDirectory() as tmp:
        _patch(Package(dist), spec, Path(tmp))


def _patch(pkg: Package, spec: PatchSpec, tmp: Path) -> None:
    if not pkg.files:
        raise EmptyDistributionRecordsError(pkg.name)

    _copy_files(pkg.path, tmp, pkg.files)

    try:
        result = _PatchApplier(tmp, spec.strip).apply(spec.patches)
    except InvalidPatchError as err:
        raise PackagePatchError(pkg.name) from err

    if result == _PatchApplier.Result.ALREADY_PATCHED:
        raise AlreadyPatchedError(pkg.name)

    copied_files = _copy_tree(tmp, pkg.path)
    _TreeCleaner(pkg.path).clean(set(pkg.files) - copied_files)


def _copy_files(src: Path, dst: Path, files: Iterable[PackagePath]) -> None:
    try:
        for f in files:
            # the file may have been deleted by a previous patch.
            if not (src / f).exists():
                continue

            (dst / f).parent.mkdir(parents=True, exist_ok=True)

            shutil.copy2(src / f, dst / f)
    except OSError as ex:
        raise CopyError(src, dst) from ex


def _copy_tree(src: Path, dst: Path) -> set[Path]:
    copied_files = set()

    def copy(s: str, d: str) -> None:
        shutil.copy2(s, d)
        copied_files.add(Path(s).relative_to(src))

    try:
        shutil.copytree(src, dst, copy_function=copy, dirs_exist_ok=True)
    except OSError as ex:
        raise CopyError(src, dst) from ex

    return copied_files


class _TreeCleaner:
    def __init__(self, root: Path) -> None:
        self._root = root

    def clean(self, files: set[PackagePath]) -> None:
        """Remove files from the root, stale .pyc caches, and empty parents."""

        for f in files:
            abs_path = self._root / f

            self._clean_file(abs_path)
            self._clean_empty_parents(abs_path)
            self._clean_pycache(abs_path)

    @staticmethod
    def _clean_file(path: Path) -> None:
        try:
            path.unlink(missing_ok=True)
        except PermissionError as ex:
            raise UnlinkError(path) from ex

    def _clean_empty_parents(self, path: Path) -> None:
        """Clean up empty parent folders up to self._root."""

        for parent in path.parents:
            if parent == self._root:
                break

            try:
                if not all(
                    p.name == "__pycache__" and p.is_dir()
                    for p in parent.iterdir()  # may throw if it was removed
                ):
                    break

                shutil.rmtree(parent / "__pycache__", ignore_errors=True)
                parent.rmdir()
            except FileNotFoundError:
                continue  # already gone by previous iteration, continue
            except OSError as ex:
                raise RemoveFolderError(parent) from ex

    @staticmethod
    def _clean_pycache(path: Path) -> None:
        if path.suffix != ".py":
            return

        pycache = path.parent / "__pycache__"
        for pyc in pycache.glob(f"{path.stem}.*.pyc"):
            # make best effort to clean up cache - skip if errors.
            with suppress(OSError):
                pyc.unlink(missing_ok=True)


class _PatchApplier:
    class PatchFlag(StrEnum):
        """PatchFlag enumerates selected flags from patch(1)."""

        # keep-sorted start
        DIRECTORY = "--directory"
        FORWARD = "--forward"
        INPUT = "--input"
        REMOVE_EMPTY_FILES = "--remove-empty-files"
        REVERSE = "--reverse"
        SILENT = "--silent"
        STRIP = "--strip"
        V = "-V"
        # keep-sorted end

    class PatchBackupMethod(StrEnum):
        NONE = auto()

    class Result(IntEnum):
        PATCHED = auto()
        ALREADY_PATCHED = auto()

    def __init__(self, root: Path, strip: int) -> None:
        self._root = root
        self._strip = strip

        cmd = shutil.which("patch")
        if cmd is None:
            raise PatchNotFoundError()
        self._patch_cmd = cmd

    def apply(self, patches: tuple[Path, ...]) -> Result:
        """Apply patch files to the root folder.

        It attempts to patch root with files. If the first patch fails, it
        attempts to detect whether root is already patched by replaying the
        reversed patch files in reverse order.

        In any case, patch() leaves root in modified state unless the first
        patch is invalid.

        Returns:
            Result to indicate the patch status.

        Raises:
            InvalidPatchError: if any of the patches fail and the root is not
              already patched with files.
        """

        for i, patch in enumerate(patches):
            cmd = self._patch_command(patch) + (
                _PatchApplier.PatchFlag.FORWARD,
            )

            proc = subprocess.run(
                tuple(str(a) for a in cmd),
                text=True,
                stdin=subprocess.DEVNULL,
                stderr=subprocess.STDOUT,
                stdout=subprocess.PIPE,
                start_new_session=True,  # suppress TTY
            )
            if proc.returncode != os.EX_OK:
                if i == 0 and self._is_patched(patches):
                    return self.Result.ALREADY_PATCHED

                raise InvalidPatchError(patch, proc.stdout.rstrip())

        return self.Result.PATCHED

    def _is_patched(self, patches: tuple[Path, ...]) -> bool:
        """Apply reversed patches in the reverse order.

        The function leaves the root in patched state.

        Args:
            files: a list of patch files.

        Returns:
            True: the reversed patches apply cleanly.
            False: otherwise
        """

        for patch in reversed(patches):
            cmd = self._patch_command(patch) + (
                _PatchApplier.PatchFlag.REVERSE,
            )

            proc = subprocess.run(
                tuple(str(a) for a in cmd),
                text=True,
                stdin=subprocess.DEVNULL,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
                start_new_session=True,  # suppress TTY
            )
            if proc.returncode != os.EX_OK:
                return False

        return True

    def _patch_command(self, patch: Path) -> tuple:
        return (
            self._patch_cmd,
            # keep-sorted start group_start_regex=['^\s+_PatchApplier.*Flag']
            _PatchApplier.PatchFlag.DIRECTORY,
            self._root,
            _PatchApplier.PatchFlag.INPUT,
            patch,
            _PatchApplier.PatchFlag.REMOVE_EMPTY_FILES,
            _PatchApplier.PatchFlag.SILENT,
            _PatchApplier.PatchFlag.STRIP,
            self._strip,
            _PatchApplier.PatchFlag.V,
            _PatchApplier.PatchBackupMethod.NONE,
            # keep-sorted end
        )
