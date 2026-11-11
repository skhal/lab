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
    Path(__file__).parents[4] / "roles/jail_skel/library/user_migrate.py"
)


@pytest.fixture
def user_migrate():
    """Loads user_migrate.py with a fake ansible.module_utils.basic.

    See zfs_create_test.py for why the fake is needed.
    """
    fake_basic = types.SimpleNamespace(AnsibleModule=object)
    with mock.patch.dict(
        sys.modules, {"ansible.module_utils.basic": fake_basic}
    ):
        spec = importlib.util.spec_from_file_location(
            "user_migrate", _MODULE_PATH
        )
        assert spec is not None
        assert spec.loader is not None
        mod = importlib.util.module_from_spec(spec)
        # dataclasses look the module up by name while it executes.
        sys.modules["user_migrate"] = mod
        spec.loader.exec_module(mod)
        yield mod


@pytest.mark.parametrize(
    "users, want",
    [
        pytest.param(["foo"], "users 'foo' not found", id="one"),
        pytest.param(["foo", "bar"], "users 'foo,bar' not found", id="many"),
    ],
)
def test_user_not_found_error(user_migrate, users, want):
    err = user_migrate.UserNotFoundError(*users)

    assert str(err) == want


def test_pw_rootdir_defaults_to_none(user_migrate):
    commands: list[Sequence[str]] = []

    def run_command(cmd: Sequence[str]) -> tuple[int, str, str]:
        commands.append(cmd)
        return os.EX_OK, "", ""

    module = AnsibleModuleProxy(run_command=run_command, check_mode=False)
    pw = user_migrate.PW(module)
    user = user_migrate.User(
        name="foo", uid=1001, gid=1001, gecos="", home="/home/foo", shell="sh"
    )

    err = pw.useradd(user, [])

    assert err is None
    assert user_migrate.PW.Flag.Rootdir not in commands[0]


def test_pw_useradd_error_reports_groups(user_migrate):
    def run_command(cmd: Sequence[str]) -> tuple[int, str, str]:
        return 1, "", "boom"

    module = AnsibleModuleProxy(run_command=run_command, check_mode=False)
    pw = user_migrate.PW(module)
    user = user_migrate.User(
        name="foo", uid=1001, gid=1001, gecos="", home="/home/foo", shell="sh"
    )

    err = pw.useradd(user, [10, 20])

    assert err is not None
    assert "groups:[10, 20]" in str(err)
