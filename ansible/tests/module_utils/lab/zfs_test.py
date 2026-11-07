# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

from dataclasses import KW_ONLY, dataclass, field
from typing import Callable, Optional, Sequence, override

import pytest

from ansible.module_utils.lab.errors import Error
from ansible.module_utils.lab.module import AnsibleModuleProxy
from ansible.module_utils.lab.zfs import ZFS


class FakeAnsibleModule(AnsibleModuleProxy):
    def __init__(
        self,
        *,
        run_command_return: tuple[int, str, str],
        check_mode: bool,
    ):
        super().__init__(
            run_command=lambda _: run_command_return, check_mode=check_mode
        )
        self._cmd = None
        self._params = dict(
            run_command_return=run_command_return, check_mode=check_mode
        )

    @property
    def cmd(self):
        return self._cmd

    @override
    def run_command(self, cmd: Sequence[str]) -> tuple[int, str, str]:
        if self._cmd is None:
            self._cmd = [cmd]
        else:
            self._cmd.append(cmd)
        return super().run_command(cmd)

    def __repr__(self):
        return f"FakeAnsibleModule(run_command_return={self._params['run_command_return']}, check_mode={self._params['check_mode']})"


@dataclass(frozen=True)
class ModuleParams:
    __test__ = False

    _: KW_ONLY
    run_command_return: tuple[int, str, str] = field(default=(0, "", ""))
    check_mode: bool = field(default=False)


@pytest.fixture
def module(request) -> FakeAnsibleModule:
    test_case: ModuleParams = request.param
    return FakeAnsibleModule(
        run_command_return=test_case.run_command_return,
        check_mode=test_case.check_mode,
    )


class TestZFS_Get:
    @dataclass(frozen=True)
    class Args:
        __test__ = False

        _: KW_ONLY
        name: str
        props: list[ZFS.Property]

    @dataclass(frozen=True)
    class Want:
        __test__ = False

        _: KW_ONLY
        err: bool = False
        cmd: list[str] = field(default_factory=list)
        props: list[str] = field(default_factory=list)

    _PARAMS = (
        pytest.param(
            ModuleParams(run_command_return=(0, "zroot/test-dset-name", "")),
            Args(
                name="zroot/test-dset",
                props=[ZFS.Property.NAME],
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.GET,
                    ZFS.GetFlag.NO_HEADER,
                    ZFS.GetFlag.PARSABLE_VALUE,
                    ZFS.GetFlag.COLUMN,
                    "value",
                    ZFS.Property.NAME,
                    "zroot/test-dset",
                ],
                props=[
                    "zroot/test-dset-name",
                ],
            ),
            id="dataset exists, get name",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(0, "zroot/test-dset-name", ""),
                check_mode=True,
            ),
            Args(
                name="zroot/test-dset",
                props=[ZFS.Property.NAME],
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.GET,
                    ZFS.GetFlag.NO_HEADER,
                    ZFS.GetFlag.PARSABLE_VALUE,
                    ZFS.GetFlag.COLUMN,
                    "value",
                    ZFS.Property.NAME,
                    "zroot/test-dset",
                ],
                props=[
                    "zroot/test-dset-name",
                ],
            ),
            id="check mode: dataset exists, get name",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(1, "", "zroot/test-dset does not exist"),
            ),
            Args(
                name="zroot/test-dset",
                props=[ZFS.Property.NAME],
            ),
            Want(
                err=True,
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.GET,
                    ZFS.GetFlag.NO_HEADER,
                    ZFS.GetFlag.PARSABLE_VALUE,
                    ZFS.GetFlag.COLUMN,
                    "value",
                    ZFS.Property.NAME,
                    "zroot/test-dset",
                ],
            ),
            id="dataset does not exist, get name",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(1, "", "zroot/test-dset does not exist"),
                check_mode=True,
            ),
            Args(
                name="zroot/test-dset",
                props=[ZFS.Property.NAME],
            ),
            Want(
                err=True,
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.GET,
                    ZFS.GetFlag.NO_HEADER,
                    ZFS.GetFlag.PARSABLE_VALUE,
                    ZFS.GetFlag.COLUMN,
                    "value",
                    ZFS.Property.NAME,
                    "zroot/test-dset",
                ],
            ),
            id="check mode: dataset does not exist, get name",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(
                    0,
                    "zroot/test-dset-name\nzroot/test-dset-mountpoint",
                    "",
                ),
            ),
            Args(
                name="zroot/test-dset",
                props=[ZFS.Property.NAME, ZFS.Property.MOUNTPOINT],
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.GET,
                    ZFS.GetFlag.NO_HEADER,
                    ZFS.GetFlag.PARSABLE_VALUE,
                    ZFS.GetFlag.COLUMN,
                    "value",
                    ",".join((ZFS.Property.NAME, ZFS.Property.MOUNTPOINT)),
                    "zroot/test-dset",
                ],
                props=[
                    "zroot/test-dset-name",
                    "zroot/test-dset-mountpoint",
                ],
            ),
            id="dataset exists, get multiple properties",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(
                    0,
                    "zroot/test-dset-name\nzroot/test-dset-mountpoint",
                    "",
                ),
                check_mode=True,
            ),
            Args(
                name="zroot/test-dset",
                props=[ZFS.Property.NAME, ZFS.Property.MOUNTPOINT],
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.GET,
                    ZFS.GetFlag.NO_HEADER,
                    ZFS.GetFlag.PARSABLE_VALUE,
                    ZFS.GetFlag.COLUMN,
                    "value",
                    ",".join((ZFS.Property.NAME, ZFS.Property.MOUNTPOINT)),
                    "zroot/test-dset",
                ],
                props=[
                    "zroot/test-dset-name",
                    "zroot/test-dset-mountpoint",
                ],
            ),
            id="check mode: dataset exists, get multiple properties",
        ),
    )

    @pytest.mark.parametrize(
        "module,args,want",
        _PARAMS,
        indirect=["module"],
    )
    @staticmethod
    def test(module: FakeAnsibleModule, args: Args, want: Want) -> None:
        zfs = ZFS(module)

        props, err = zfs.get(args.name, *args.props)

        assert module.cmd == [want.cmd]
        if want.err:
            assert err is not None
        else:
            assert err is None
        assert props == want.props


class TestZFS_Clone:
    @dataclass(frozen=True)
    class Args:
        __test__ = False

        _: KW_ONLY
        snapshot: str
        dset: str
        props: dict[ZFS.Property, str] = field(default_factory=dict)

    @dataclass(frozen=True)
    class Want:
        __test__ = False

        _: KW_ONLY
        err: bool = False
        cmd: Optional[list[str]] = field(default=None)

    _PARAMS = (
        pytest.param(
            ModuleParams(),
            Args(
                snapshot="test-origin@snapshot",
                dset="test-target",
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.CLONE,
                    "test-origin@snapshot",
                    "test-target",
                ],
            ),
            id="target dataset does not exist",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(1, "", "dataset test-target exists"),
            ),
            Args(
                snapshot="test-origin@snapshot",
                dset="test-target",
            ),
            Want(
                err=True,
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.CLONE,
                    "test-origin@snapshot",
                    "test-target",
                ],
            ),
            id="target dataset exists",
        ),
        pytest.param(
            ModuleParams(),
            Args(
                snapshot="test-origin@snapshot",
                dset="test-target",
                props={ZFS.Property.CAN_MOUNT: "on"},
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.CLONE,
                    ZFS.CloneFlag.OPTION,
                    "canmount=on",
                    "test-origin@snapshot",
                    "test-target",
                ],
            ),
            id="pass one property",
        ),
        pytest.param(
            ModuleParams(),
            Args(
                snapshot="test-origin@snapshot",
                dset="test-target",
                props={
                    ZFS.Property.CAN_MOUNT: "on",
                    ZFS.Property.READ_ONLY: "off",
                },
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    ZFS.Command.CLONE,
                    ZFS.CloneFlag.OPTION,
                    "canmount=on",
                    ZFS.CloneFlag.OPTION,
                    "readonly=off",
                    "test-origin@snapshot",
                    "test-target",
                ],
            ),
            id="pass multiple properties",
        ),
        pytest.param(
            ModuleParams(check_mode=True),
            Args(
                snapshot="test-origin@snapshot",
                dset="test-target",
            ),
            Want(),
            id="check mode: target dataset does not exist",
        ),
    )

    @pytest.mark.parametrize("module,args,want", _PARAMS, indirect=["module"])
    @staticmethod
    def test(module: FakeAnsibleModule, args: Args, want: Want) -> None:
        zfs = ZFS(module)

        err = zfs.clone(args.snapshot, args.dset, args.props)

        assert module.cmd == (None if want.cmd is None else [want.cmd])
        if want.err:
            assert err is not None
        else:
            assert err is None


class TestZFS_Create:
    @dataclass(frozen=True)
    class Args:
        __test__ = False

        _: KW_ONLY
        dset: str
        props: dict[ZFS.Property, str] = field(default_factory=dict)
        parent_props: dict[ZFS.Property, str] = field(default_factory=dict)

    @dataclass(frozen=True)
    class Want:
        __test__ = False

        _: KW_ONLY
        err: bool = False
        cmd: list[list[str]] = field(default_factory=list)

    _PARAMS = (
        pytest.param(
            ModuleParams(),
            Args(
                dset="zroot/test-target",
            ),
            Want(
                cmd=[
                    [
                        "/sbin/zfs",
                        ZFS.Command.CREATE,
                        "zroot/test-target",
                    ]
                ],
            ),
            id="target dataset does not exist",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(1, "", "dataset zroot/test-target exists"),
            ),
            Args(
                dset="zroot/test-target",
            ),
            Want(
                err=True,
                cmd=[
                    [
                        "/sbin/zfs",
                        ZFS.Command.CREATE,
                        "zroot/test-target",
                    ]
                ],
            ),
            id="target dataset exists",
        ),
        pytest.param(
            ModuleParams(),
            Args(
                dset="zroot/test-target",
                props={ZFS.Property.CAN_MOUNT: "off"},
            ),
            Want(
                cmd=[
                    [
                        "/sbin/zfs",
                        ZFS.Command.CREATE,
                        ZFS.CreateFlag.OPTION,
                        "canmount=off",
                        "zroot/test-target",
                    ]
                ],
            ),
            id="pass one property",
        ),
        pytest.param(
            ModuleParams(),
            Args(
                dset="zroot/test-target",
                props={
                    ZFS.Property.CAN_MOUNT: "off",
                    ZFS.Property.READ_ONLY: "on",
                },
            ),
            Want(
                cmd=[
                    [
                        "/sbin/zfs",
                        ZFS.Command.CREATE,
                        ZFS.CreateFlag.OPTION,
                        "canmount=off",
                        ZFS.CreateFlag.OPTION,
                        "readonly=on",
                        "zroot/test-target",
                    ]
                ],
            ),
            id="pass multiple properties",
        ),
    )

    @pytest.mark.parametrize("module,args,want", _PARAMS, indirect=["module"])
    @staticmethod
    def test(module: FakeAnsibleModule, args: Args, want: Want) -> None:
        zfs = ZFS(module)

        err = zfs.create(args.dset, args.props, args.parent_props)

        assert module.cmd == want.cmd
        if want.err:
            assert err is not None
        else:
            assert err is None


class TestZFS_List:
    @dataclass(frozen=True)
    class Args:
        __test__ = False

        _: KW_ONLY
        dset: str
        fstype: ZFS.FSType

    @dataclass(frozen=True)
    class Want:
        __test__ = False

        _: KW_ONLY
        err: bool = False
        cmd: list[str] = field(default_factory=list)
        out: list[str] = field(default_factory=list)

    _PARAMS = (
        pytest.param(
            ModuleParams(run_command_return=(0, "zroot/test-dset@one", "")),
            Args(
                dset="zroot/test-dset",
                fstype=ZFS.FSType.SNAPSHOT,
            ),
            Want(
                cmd=[
                    "/sbin/zfs",
                    "list",
                    "-H",
                    "-o",
                    "name",
                    "-t",
                    "snapshot",
                    "zroot/test-dset",
                ],
                out=[
                    "zroot/test-dset@one",
                ],
            ),
            id="dataset exists",
        ),
        pytest.param(
            ModuleParams(
                run_command_return=(
                    1,
                    "",
                    "cannot open 'zroot/test-dset': dataset does not exist",
                )
            ),
            Args(
                dset="zroot/test-dset",
                fstype=ZFS.FSType.SNAPSHOT,
            ),
            Want(
                err=True,
                cmd=[
                    "/sbin/zfs",
                    "list",
                    "-H",
                    "-o",
                    "name",
                    "-t",
                    "snapshot",
                    "zroot/test-dset",
                ],
            ),
            id="dataset does not exist",
        ),
    )

    @pytest.mark.parametrize("module,args,want", _PARAMS, indirect=["module"])
    @staticmethod
    def test(module: FakeAnsibleModule, args: Args, want: Want) -> None:
        zfs = ZFS(module)

        out, err = zfs.list(args.dset, args.fstype)

        assert module.cmd == [want.cmd]
        if want.err:
            assert err is not None
        else:
            assert err is None
        assert out == want.out
