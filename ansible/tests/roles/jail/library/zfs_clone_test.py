# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

import importlib.util
import os
import sys
import types
from pathlib import Path
from typing import Sequence
from unittest import mock

import pytest

from ansible.module_utils.lab.module import AnsibleModuleProxy

_MODULE_PATH = Path(__file__).parents[4] / "roles/jail/library/zfs_clone.py"


@pytest.fixture
def zfs_clone():
    """Loads zfs_clone.py with a fake ansible.module_utils.basic.

    See zfs_create_test.py for why ansible.module_utils.basic is faked.
    """
    fake_basic = types.SimpleNamespace(AnsibleModule=object)
    with mock.patch.dict(
        sys.modules, {"ansible.module_utils.basic": fake_basic}
    ):
        spec = importlib.util.spec_from_file_location("zfs_clone", _MODULE_PATH)
        assert spec is not None
        assert spec.loader is not None
        mod = importlib.util.module_from_spec(spec)
        # dataclasses resolves string annotations through sys.modules.
        sys.modules[spec.name] = mod
        spec.loader.exec_module(mod)
        yield mod


class FakeZFS:
    """Fakes `zfs get`, `zfs list` and `zfs clone` over datasets.

    datasets maps a dataset name to its (origin, mountpoint). snapshots maps
    a dataset name to its snapshots.
    """

    def __init__(
        self,
        datasets: dict[str, tuple[str, str]] | None = None,
        snapshots: dict[str, list[str]] | None = None,
    ) -> None:
        self.datasets = datasets or {}
        self.snapshots = snapshots or {}
        self.cloned: list[tuple[str, str]] = []

    def run_command(self, cmd: Sequence[str]) -> tuple[int, str, str]:
        if "get" in cmd:
            dset = self.datasets.get(cmd[-1])
            if dset is None:
                return 1, "", "dataset does not exist"
            return os.EX_OK, "\n".join(dset) + "\n", ""

        if "list" in cmd:
            snaps = self.snapshots.get(cmd[-1], [])
            return os.EX_OK, "".join(s + "\n" for s in snaps), ""

        if "clone" in cmd:
            origin, name = cmd[-2], cmd[-1]
            if name in self.datasets:
                return 1, "", "dataset already exists"
            self.cloned.append((origin, name))
            self.datasets[name] = (origin, "/" + name)
            return os.EX_OK, "", ""

        return 1, "", "unexpected command"

    def module(self, *, check_mode: bool = False) -> AnsibleModuleProxy:
        return AnsibleModuleProxy(
            run_command=self.run_command, check_mode=check_mode
        )


def test_run_clone_new_dataset(zfs_clone):
    fake = FakeZFS(snapshots={"zroot/src": ["zroot/src@a", "zroot/src@b"]})
    params = dict(name="zroot/dst", origin="zroot/src@a", props={})

    got, err = zfs_clone.run(fake.module(), params)

    assert err is None
    assert fake.cloned == [("zroot/src@a", "zroot/dst")]
    assert got["origin"] == "zroot/src@a"
    assert got["changed"]


def test_run_clone_latest_snapshot(zfs_clone):
    fake = FakeZFS(snapshots={"zroot/src": ["zroot/src@a", "zroot/src@b"]})
    params = dict(name="zroot/dst", origin="zroot/src", props={})

    got, err = zfs_clone.run(fake.module(), params)

    assert err is None
    assert fake.cloned == [("zroot/src@b", "zroot/dst")]
    assert got["origin"] == "zroot/src@b"
    assert got["changed"]


@pytest.mark.parametrize("origin", ["zroot/src@a", "zroot/src"])
def test_run_existing_dataset_is_idempotent(zfs_clone, origin):
    fake = FakeZFS(
        datasets={"zroot/dst": ("zroot/src@a", "/dst")},
        snapshots={"zroot/src": ["zroot/src@a"]},
    )
    params = dict(name="zroot/dst", origin=origin, props={})

    got, err = zfs_clone.run(fake.module(), params)

    assert err is None
    assert fake.cloned == []
    assert got["name"] == "zroot/dst"
    assert got["origin"] == "zroot/src@a"
    assert got["mountpoint"] == "/dst"
    assert not got["changed"]


def test_run_existing_dataset_origin_mismatch(zfs_clone):
    fake = FakeZFS(datasets={"zroot/dst": ("zroot/src@a", "/dst")})
    params = dict(name="zroot/dst", origin="zroot/src@b", props={})

    _, err = zfs_clone.run(fake.module(), params)

    assert err is not None
    assert "origin 'zroot/src@a'" in str(err)
    assert "requested origin 'zroot/src@b'" in str(err)
    assert fake.cloned == []
