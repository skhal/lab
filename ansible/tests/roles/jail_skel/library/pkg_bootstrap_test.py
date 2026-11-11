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
from unittest import mock

import pytest

from ansible.module_utils.lab.module import AnsibleModuleProxy

_MODULE_PATH = (
    Path(__file__).parents[4] / "roles/jail_skel/library/pkg_bootstrap.py"
)


@pytest.fixture
def pkg_bootstrap():
    """Loads pkg_bootstrap.py with a fake ansible.module_utils.basic.

    The module is registered in sys.modules so that pyfakefs patches its os and
    shutil references.
    """
    fake_basic = types.SimpleNamespace(AnsibleModule=object)
    with mock.patch.dict(
        sys.modules, {"ansible.module_utils.basic": fake_basic}
    ):
        spec = importlib.util.spec_from_file_location(
            "pkg_bootstrap", _MODULE_PATH
        )
        assert spec is not None
        assert spec.loader is not None
        mod = importlib.util.module_from_spec(spec)
        sys.modules["pkg_bootstrap"] = mod
        spec.loader.exec_module(mod)
        yield mod


_KEYS = "/usr/share/keys"
_ROOT = "/jail"


def new_pkg(pkg_bootstrap, *, check_mode: bool):
    module = AnsibleModuleProxy(
        run_command=lambda cmd: (os.EX_OK, "", ""), check_mode=check_mode
    )
    return pkg_bootstrap.PKG(module)


@pytest.fixture
def keys(pkg_bootstrap, fs):
    """Creates source keys in the fake filesystem.

    The pkg_bootstrap fixture must come first so that the module is loaded
    from the real filesystem, before pyfakefs replaces it.
    """
    fs.create_file(f"{_KEYS}/pkg/trusted/a", contents="A")
    fs.create_file(f"{_KEYS}/pkg/revoked/b", contents="B")
    fs.create_dir(_ROOT)
    return _KEYS


def dest(path: str) -> Path:
    return Path(_ROOT + _KEYS) / path


def test_fresh_install(pkg_bootstrap, keys):
    pkg = new_pkg(pkg_bootstrap, check_mode=False)

    err = pkg.bootstrap(_ROOT)

    assert err is None, str(err)
    assert pkg.changed
    assert len(pkg.diff) == 2
    assert dest("pkg/trusted/a").read_text() == "A"
    assert dest("pkg/revoked/b").read_text() == "B"


def test_existing_dest_dir(pkg_bootstrap, keys, fs):
    fs.create_dir(_ROOT + _KEYS)
    pkg = new_pkg(pkg_bootstrap, check_mode=False)

    err = pkg.bootstrap(_ROOT)

    assert err is None, str(err)
    assert dest("pkg/trusted/a").read_text() == "A"


def test_idempotent(pkg_bootstrap, keys):
    assert new_pkg(pkg_bootstrap, check_mode=False).bootstrap(_ROOT) is None
    pkg = new_pkg(pkg_bootstrap, check_mode=False)

    err = pkg.bootstrap(_ROOT)

    assert err is None, str(err)
    assert not pkg.changed
    assert pkg.diff == []


def test_outdated_file_is_replaced(pkg_bootstrap, keys, fs):
    fs.create_file(str(dest("pkg/trusted/a")), contents="old")
    pkg = new_pkg(pkg_bootstrap, check_mode=False)

    err = pkg.bootstrap(_ROOT)

    assert err is None, str(err)
    assert dest("pkg/trusted/a").read_text() == "A"
    assert {d["before"]["state"] for d in pkg.diff} == {"outdated", "absent"}


def test_check_mode(pkg_bootstrap, keys):
    pkg = new_pkg(pkg_bootstrap, check_mode=True)

    err = pkg.bootstrap(_ROOT)

    assert err is None, str(err)
    assert pkg.changed
    assert len(pkg.diff) == 2
    assert not Path(_ROOT + _KEYS).exists()


def test_error_propagates(pkg_bootstrap, keys):
    pkg = new_pkg(pkg_bootstrap, check_mode=False)

    with mock.patch.object(
        pkg_bootstrap.shutil, "copy2", side_effect=OSError("boom")
    ):
        err = pkg.bootstrap(_ROOT)

    assert err is not None
    assert "boom" in str(err)


@pytest.mark.parametrize("check_mode", [False, True])
def test_missing_source_fails(pkg_bootstrap, fs, check_mode):
    fs.create_dir(_ROOT)
    pkg = new_pkg(pkg_bootstrap, check_mode=check_mode)

    err = pkg.bootstrap(_ROOT)

    assert err is not None
    assert _KEYS in str(err)
    assert not pkg.changed
    assert os.listdir(_ROOT) == []
