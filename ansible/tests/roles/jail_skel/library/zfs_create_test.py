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

_MODULE_PATH = (
    Path(__file__).parents[4] / "roles/jail_skel/library/zfs_create.py"
)


@pytest.fixture
def zfs_create():
    """Loads zfs_create.py with a fake ansible.module_utils.basic.

    Under importlib import mode the local ansible/ folder shadows the installed
    ansible package, so ansible.module_utils.basic is unavailable. The module
    under test only needs the AnsibleModule name at import time. sys.modules is
    restored when the test ends.
    """
    fake_basic = types.SimpleNamespace(AnsibleModule=object)
    with mock.patch.dict(
        sys.modules, {"ansible.module_utils.basic": fake_basic}
    ):
        spec = importlib.util.spec_from_file_location(
            "zfs_create", _MODULE_PATH
        )
        assert spec is not None
        assert spec.loader is not None
        mod = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(mod)
        yield mod


def new_module(
    mountpoints: dict[str, str], *, check_mode: bool
) -> AnsibleModuleProxy:
    """Fakes `zfs get` and `zfs create` over a set of existing datasets."""

    def run_command(cmd: Sequence[str]) -> tuple[int, str, str]:
        if "get" in cmd:
            mountpoint = mountpoints.get(cmd[-1])
            if mountpoint is None:
                return 1, "", "dataset does not exist"
            return os.EX_OK, mountpoint + "\n", ""
        return os.EX_OK, "", ""

    return AnsibleModuleProxy(run_command=run_command, check_mode=check_mode)


@pytest.mark.parametrize(
    "mountpoints, dset, want_mountpoint",
    [
        pytest.param(
            {"zroot/foo": "/example/foo"},
            "zroot/foo",
            "/example/foo",
            id="existing",
        ),
        pytest.param(
            {"zroot/foo": "/example/foo"},
            "zroot/foo/bar/baz",
            "/example/foo/bar/baz",
            id="missing_children",
        ),
    ],
)
def test_get_mountpoint_in_check_mode(
    zfs_create, mountpoints, dset, want_mountpoint
):
    module = new_module(mountpoints, check_mode=True)
    zfs = zfs_create.ZFS(module)

    got, err = zfs_create.get_mountpoint_in_check_mode(zfs, dset)

    assert err is None
    assert got == want_mountpoint


def test_get_mountpoint_in_check_mode_no_parent(zfs_create):
    module = new_module({}, check_mode=True)
    zfs = zfs_create.ZFS(module)

    _, err = zfs_create.get_mountpoint_in_check_mode(zfs, "zroot/foo")

    assert err is not None


def test_run_check_mode_new_dataset(zfs_create):
    module = new_module({"zroot/foo": "/example/foo"}, check_mode=True)
    params = dict(name="zroot/foo/bar", props={}, parent_props={})

    got, err = zfs_create.run(module, params)

    assert err is None
    assert got["mountpoint"] == "/example/foo/bar"
    assert got["changed"]
