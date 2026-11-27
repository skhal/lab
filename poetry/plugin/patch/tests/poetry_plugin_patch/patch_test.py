# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


import difflib
from contextlib import AbstractContextManager, nullcontext
from dataclasses import dataclass, field
from importlib.metadata import Distribution, PathDistribution
from pathlib import Path
from typing import Any

import pytest

from poetry_plugin_patch.error import (
    AlreadyPatchedError,
    EmptyDistributionRecordsError,
    PackagePatchError,
)
from poetry_plugin_patch.patch import patch
from poetry_plugin_patch.spec import PatchSpec


@dataclass(frozen=True)
class DistributionParams:
    name: str
    files: dict[str, str]
    # unlisted files are present in the distribution but are not in RECORD, e.g.
    # Python cache files .pyc.
    unlisted: dict[str, str] = field(default_factory=dict)
    version: str = "0.0.0"


@pytest.fixture(name="dist")
def parse_distspec(
    request: pytest.FixtureRequest, tmp_path: Path
) -> Distribution:
    params: DistributionParams = request.param

    site = tmp_path / "site"

    info = site / f"{params.name}-{params.version}.dist-info"
    info.mkdir(parents=True)

    records = []

    def _write(site: Path, rel_path: Path, data: str):
        abs_path = site / rel_path
        abs_path.parent.mkdir(parents=True, exist_ok=True)
        abs_path.write_text(data)

    for rel_path, data in params.files.items():
        _write(site, Path(rel_path), data)
        records.append(f"{rel_path},,")

    for rel_path, data in params.unlisted.items():
        _write(site, Path(rel_path), data)

    meta = info / "METADATA"
    records.append(f"{Path(info.name) / meta.name},,")

    meta.write_text(
        f"""\
Metadata-Version: 2.1
Name: {params.name}
Version: {params.version}
"""
    )

    rec = info / "RECORD"
    records.append(f"{Path(info.name) / rec.name},,")

    rec.write_text(f"""\
{"\n".join(records)}
""")

    return PathDistribution(info)


@dataclass(frozen=True)
class PatchParams:
    file: str


@dataclass(frozen=True)
class ModifyPatchParams(PatchParams):
    old: str
    new: str


@dataclass(frozen=True)
class DeletePatchParams(PatchParams):
    old: str | None = None  # None reads the file contents from the distribution


@dataclass(frozen=True)
class PatchSpecParams:
    package: str
    patches: dict[str, PatchParams]
    strip: int


@pytest.fixture(name="spec")
def parse_patchspec(
    request: pytest.FixtureRequest, tmp_path: Path, dist: Distribution
) -> PatchSpec:
    params: PatchSpecParams = request.param

    dist_site = dist.locate_file("")

    patches = []
    for rel_patch, patch_params in params.patches.items():
        match patch_params:
            case ModifyPatchParams(file=file, old=old, new=new):
                diff = difflib.unified_diff(
                    old.splitlines(keepends=True),
                    new.splitlines(keepends=True),
                    fromfile=f"{Path('a') / file}",
                    tofile=f"{Path('b') / file}",
                )
            case DeletePatchParams(file=file, old=old):
                if old is None:
                    old = (dist_site / file).read_text()
                diff = difflib.unified_diff(
                    old.splitlines(keepends=True),
                    [],
                    fromfile=f"{Path('a') / file}",
                    tofile="/dev/null",
                )
            case _:
                raise TypeError(f"unsupported patch params: {patch_params!r}")

        patch = tmp_path / rel_patch
        patch.write_text("".join(diff))

        patches.append(rel_patch)

    spec = dict(
        package=params.package,
        patches=patches,
        strip=params.strip,
    )

    return PatchSpec.parse(spec, tmp_path)


def snapshot(dist: Distribution) -> dict[str, str]:
    path = dist.locate_file("")
    assert isinstance(path, Path)

    tree = {}
    for p in path.rglob("*"):
        rel = p.relative_to(path)
        if rel.parts[0].endswith(".dist-info"):
            continue
        if p.is_dir():
            tree[f"{rel}/"] = ""
        else:
            tree[str(rel)] = p.read_text()
    return tree


@dataclass(frozen=True)
class Want:
    tree: dict[str, str] | None = None  # None means want no changes
    expect: AbstractContextManager[Any] = field(default_factory=nullcontext)


_TESTS = (
    pytest.param(
        DistributionParams(name="test", files={}),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": ModifyPatchParams(
                    file="a.txt",
                    old="text\n",
                    new="new text\n",
                ),
            },
            strip=1,
        ),
        Want(
            expect=pytest.raises(EmptyDistributionRecordsError),
        ),
        id="empty distribution",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "b.patch": ModifyPatchParams(
                    file="b.txt",
                    old="text\n",
                    new="new text\n",
                ),
            },
            strip=1,
        ),
        Want(
            expect=pytest.raises(PackagePatchError),
        ),
        id="patch nonexistent file",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": ModifyPatchParams(
                    file="a.txt",
                    old="non existent text\n",
                    new="new text\n",
                ),
            },
            strip=1,
        ),
        Want(
            expect=pytest.raises(PackagePatchError),
        ),
        id="patch does not apply",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": ModifyPatchParams(
                    file="a.txt",
                    old="text\n",
                    new="new text\n",
                ),
            },
            strip=1,
        ),
        Want(
            tree={
                "a.txt": "new text\n",
            },
        ),
        id="patch applies",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "new text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": ModifyPatchParams(
                    file="a.txt",
                    old="text\n",
                    new="new text\n",
                ),
            },
            strip=1,
        ),
        Want(
            expect=pytest.raises(AlreadyPatchedError),
        ),
        id="already applied patch",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "one.patch": ModifyPatchParams(
                    file="a.txt",
                    old="text\n",
                    new="new text\n",
                ),
                "two.patch": ModifyPatchParams(
                    file="a.txt",
                    old="non existent text\n",
                    new="should not be present text\n",
                ),
            },
            strip=1,
        ),
        Want(
            expect=pytest.raises(PackagePatchError),
        ),
        id="second patch fails",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(file="a.txt", old="other text\n"),
            },
            strip=1,
        ),
        Want(
            expect=pytest.raises(PackagePatchError),
        ),
        marks=pytest.mark.xfail(
            strict=True,
            reason="_is_patched() treats a mismatched reverse del as patched",
        ),
        id="invalid delete patch",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.py": "text\n",
                "b.py": "other text\n",
            },
            unlisted={
                "__pycache__/a.cpython-315.pyc": "should be removed\n",
                "__pycache__/b.cpython-315.pyc": "should remain\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": ModifyPatchParams(
                    file="a.py",
                    old="text\n",
                    new="new text\n",
                ),
            },
            strip=1,
        ),
        Want(
            tree={
                "a.py": "new text\n",
                "b.py": "other text\n",
                "__pycache__/": "",
                "__pycache__/b.cpython-315.pyc": "should remain\n",
            }
        ),
        marks=pytest.mark.xfail(
            strict=True,
            reason="_patch() only cleans .pyc for deleted files",
        ),
        id="removes python cache of the modified file",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.py": "text\n",
                "b.py": "other text\n",
            },
            unlisted={
                "__pycache__/a.cpython-315.pyc": "should be removed\n",
                "__pycache__/b.cpython-315.pyc": "should remain\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(
                    file="a.py",
                ),
            },
            strip=1,
        ),
        Want(
            tree={
                "b.py": "other text\n",
                "__pycache__/": "",
                "__pycache__/b.cpython-315.pyc": "should remain\n",
            }
        ),
        id="removes python cache of the deleted file",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(
                    file="a.txt",
                ),
            },
            strip=1,
        ),
        Want(tree={}),
        id="clean after delete",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "sub/a.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(
                    file="sub/a.txt",
                ),
            },
            strip=1,
        ),
        Want(tree={}),
        id="clean removes empty folders",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "sub/a.txt": "text\n",
                "sub/b.txt": "another text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(
                    file="sub/a.txt",
                ),
            },
            strip=1,
        ),
        Want(
            tree={
                "sub/": "",
                "sub/b.txt": "another text\n",
            }
        ),
        id="clean keeps non empty folders",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a/b/c.txt": "text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(
                    file="a/b/c.txt",
                ),
            },
            strip=1,
        ),
        Want(tree={}),
        id="recursive clean removes empty paths",
    ),
    pytest.param(
        DistributionParams(
            name="test",
            files={
                "a/b/c.txt": "text\n",
                "a/d.txt": "another text\n",
            },
        ),
        PatchSpecParams(
            package="test",
            patches={
                "a.patch": DeletePatchParams(
                    file="a/b/c.txt",
                ),
            },
            strip=1,
        ),
        Want(
            tree={
                "a/": "",
                "a/d.txt": "another text\n",
            }
        ),
        id="recursive clean keeps non-empty intermediate paths",
    ),
)


@pytest.mark.parametrize("dist,spec,want", _TESTS, indirect=["dist", "spec"])
def test_patch(dist: Distribution, spec: PatchSpec, want: Want):
    snap_init = snapshot(dist)
    with want.expect:
        patch(dist, spec)

    assert snapshot(dist) == (snap_init if want.tree is None else want.tree)
