# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from contextlib import AbstractContextManager, nullcontext
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

import pytest
from pyfakefs.fake_filesystem import FakeFilesystem

from poetry_plugin_patch.error import InvalidSpecError
from poetry_plugin_patch.spec import PatchSpec


@pytest.fixture
def fake_fs(request: pytest.FixtureRequest, fs: FakeFilesystem) -> None:
    files: list[str] = request.param
    for f in files:
        fs.create_file(f)


@dataclass(frozen=True)
class Want:
    spec: PatchSpec | None = None
    expect: AbstractContextManager[Any] = field(default_factory=nullcontext)


@pytest.mark.parametrize(
    "spec,base,want",
    (
        pytest.param(
            {},
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="empty spec",
        ),
        pytest.param(
            dict(package="test", patches=[], test="unexpected field"),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="unknown fields",
        ),
    ),
)
def test_parse(spec: dict, base: Path, want: Want) -> None:
    with want.expect:
        got = PatchSpec.parse(spec, base)

        assert got == want.spec


@pytest.mark.usefixtures("fake_fs")
@pytest.mark.parametrize(
    "fake_fs,spec,base,want",
    (
        pytest.param(
            [],
            dict(patches=[]),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="missing",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="", patches=["a.patch"]),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="empty",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package=123, patches=["a.patch"]),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="must be a string",
        ),
    ),
    indirect=["fake_fs"],
)
def test_parse_package(spec: dict, base: Path, want: Want) -> None:
    with want.expect:
        got = PatchSpec.parse(spec, base)

        assert got == want.spec


@pytest.mark.usefixtures("fake_fs")
@pytest.mark.parametrize(
    "fake_fs,spec,base,want",
    (
        pytest.param(
            [],
            dict(package="test"),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="missing",
        ),
        pytest.param(
            [],
            dict(package="test", patches="not a list"),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="not a list",
        ),
        pytest.param(
            [],
            dict(package="test", patches=[1, 2, 3]),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="not a list of strings",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["/tmp/a.patch"]),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="absolute path",
        ),
        pytest.param(
            [],
            dict(package="test", patches=[]),
            Path("/tmp"),
            Want(
                spec=PatchSpec(package="test", patches=tuple()),
            ),
            id="empty",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["a.patch"]),
            Path("/tmp"),
            Want(
                spec=PatchSpec(
                    package="test",
                    patches=(Path("/tmp/a.patch"),),
                    strip=1,
                ),
            ),
            id="relative patch",
        ),
    ),
    indirect=["fake_fs"],
)
def test_parse_patches(spec: dict, base: Path, want: Want) -> None:
    with want.expect:
        got = PatchSpec.parse(spec, base)

        assert got == want.spec


@pytest.mark.usefixtures("fake_fs")
@pytest.mark.parametrize(
    "fake_fs,spec,base,want",
    (
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["a.patch"], strip="test"),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="not integer",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["a.patch"], strip=True),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="bool",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["a.patch"], strip=-1),
            Path("/tmp"),
            Want(
                expect=pytest.raises(InvalidSpecError),
            ),
            id="negative",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["a.patch"]),
            Path("/tmp"),
            Want(
                spec=PatchSpec(
                    package="test",
                    patches=(Path("/tmp/a.patch"),),
                    strip=1,
                ),
            ),
            id="implicit strip",
        ),
        pytest.param(
            ["/tmp/a.patch"],
            dict(package="test", patches=["a.patch"], strip=2),
            Path("/tmp"),
            Want(
                spec=PatchSpec(
                    package="test",
                    patches=(Path("/tmp/a.patch"),),
                    strip=2,
                ),
            ),
            id="explicit strip",
        ),
    ),
    indirect=["fake_fs"],
)
def test_parse_strip(spec: dict, base: Path, want: Want) -> None:
    with want.expect:
        got = PatchSpec.parse(spec, base)

        assert got == want.spec
