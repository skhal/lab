# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

import traceback
from dataclasses import dataclass
from typing import Any

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.lab.errors import Error
from ansible.module_utils.lab.module import AnsibleModuleProxy
from ansible.module_utils.lab.zfs import ZFS


class SnapshotNotFoundError(Error):
    def __init__(self, name: str) -> None:
        super().__init__(f"snapshot '{name}' not found")


class UnsupportedSnapshotError(Error):
    def __init__(self, name: str) -> None:
        super().__init__(f"unsupported snapshot '{name}'")


@dataclass(frozen=True)
class Dataset:
    name: str
    origin: str
    mountpoint: str

    @staticmethod
    def nil() -> Dataset:
        return Dataset("", "", "")


class Cloner:
    def __init__(self, module: AnsibleModuleProxy) -> None:
        self._module = module
        self._changed = False
        self._diff = []

    @property
    def changed(self):
        return self._changed

    @property
    def diff(self):
        return self._diff

    def clone(
        self, origin: str, name: str, props: dict[ZFS.Property, str]
    ) -> tuple[Dataset | None, Error | None]:
        dset = self._maybe_get_dset_origin(name)

        if "@" not in origin:
            if dset:
                return dset, None

            snapshot, err = self._get_latest_snapshot(origin)
            if err is not None:
                return None, err

            origin = snapshot

        if dset and dset.origin != origin:
            return None, Error(
                f"dataset '{name}' origin '{dset.origin}' does not match "
                "requested origin '{origin}'"
            )

        if (err := self._clone(origin, name, props)) is not None:
            return Dataset.nil(), err

        return self._get_dset_origin(name)

    def _maybe_get_dset_origin(self, name: str) -> Dataset | None:
        dset, err = self._get_dset_origin(name)
        if err != None:
            return None

        return dset

    def _get_dset_origin(self, name: str) -> tuple[Dataset, Error | None]:
        zfs = ZFS(self._module)

        props, err = zfs.get(name, zfs.Property.ORIGIN, zfs.Property.MOUNTPOINT)
        if err is not None:
            return Dataset.nil(), err

        origin, mountpoint = props

        return Dataset(name=name, origin=origin, mountpoint=mountpoint), None

    def _get_latest_snapshot(self, dset: str) -> tuple[str, Error | None]:
        zfs = ZFS(self._module)

        snapshots, err = zfs.list(dset, ZFS.FSType.SNAPSHOT)
        if err is not None:
            return "", err

        if not snapshots:
            return "", SnapshotNotFoundError("latest")

        return snapshots[-1], None

    def _clone(
        self, origin: str, dset: str, props: dict[ZFS.Property, str]
    ) -> Error | None:
        zfs = ZFS(self._module)

        if (err := zfs.clone(origin, dset, props)) is not None:
            return err

        self._changed |= zfs._changed
        self._diff.extend(zfs.diff)


def run(
    module: AnsibleModuleProxy,
    params: frozendict[str, Any],
) -> tuple[dict, Error | None]:
    dset = params["name"]
    origin = params["origin"]
    props, err = ZFS.parse_create_properties(params["props"])
    if err != None:
        return {}, err

    cloner = Cloner(module)

    dset, err = cloner.clone(origin, dset, props)
    if err is not None:
        return {}, err

    assert dset is not None

    return (
        dict(
            name=dset.name,
            origin=dset.origin,
            mountpoint=dset.mountpoint,
            changed=cloner.changed,
            diff=cloner.diff,
        ),
        None,
    )


def main() -> None:
    spec_on_off_noauto = dict(type=str, choices=["on", "off", "noauto"])
    spec_on_off = dict(type=str, choices=["on", "off"])

    module = AnsibleModule(
        argument_spec=dict(
            name=dict(type="str", required=True),
            origin=dict(type="str", required=True),
            props=dict(
                type="dict",
                default={},
                options=dict(
                    canmount=spec_on_off_noauto,
                    readonly=spec_on_off,
                ),
            ),
        ),
        supports_check_mode=True,
    )

    module_proxy = AnsibleModuleProxy(
        run_command=lambda cmd: module.run_command(cmd),
        check_mode=module.check_mode,
    )

    try:
        res, err = run(module_proxy, module.params)
        if err is not None:
            module.fail_json(msg=str(err))

        module.exit_json(**res)
    except Exception:
        module.fail_json(msg=traceback.format_exc())


if "__main__" == __name__:
    main()
