# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""Package gives access to the distribution package.

It include a Package class that is a view into a distribution. It gives a quick
access to the distribution's path and a list of files excluding .dist-info and
__pycache__.
"""

from functools import cached_property
from importlib.metadata import Distribution, PackagePath
from pathlib import Path
from typing import cast


class Package:
    """A view into a Distribution.

    Args:
        dist: the package's distribution.
    """

    def __init__(self, dist: Distribution) -> None:  # noqa: D107
        self._dist = dist

    @property
    def name(self) -> str:
        """The name of the distribution."""
        return self._dist.name

    @cached_property
    def path(self) -> Path:
        """An absolute path to the distribution."""
        return cast(Path, self._dist.locate_file(""))

    @cached_property
    def files(self) -> tuple[PackagePath, ...]:
        """A list of distribution files excluding .dist-info and __pycache__."""
        return tuple(
            f
            for f in self._dist.files or []
            if f.parts[0] != ".."
            and not f.parts[0].endswith(".dist-info")
            and "__pycache__" not in f.parts
        )

    @property
    def version(self) -> str:
        """The version of the distribution."""
        return self._dist.version
