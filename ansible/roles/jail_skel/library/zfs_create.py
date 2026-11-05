# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

import os
import traceback
from enum import StrEnum, auto

from ansible.module_utils.basic import AnsibleModule

from ansible.module_utils.lab.errors import Error
from ansible.module_utils.lab.zfs import ZFS


def get_mountpoint_in_check_mode(
    zfs: ZFS, dset: str
) -> tuple[str, Error | None]:
    """builds a mountpoint based on the first available parent dataset.

    It does not fail if the dataset does not exist. Instead, it find the
    mountpoint of the first available parent dataset to use as a base. Then
    it appends the child datasets to that base.

    For example:

        % zfs get -H -o value mountpoint zroot/foo
        /example/foo

        % zfs list zroot/foo/bar/baz
        cannot open 'zroot/yo': dataset does not exist

        >>> get_mountpoint_in_check_mode(zfs, "zroot/foo/bar/baz")
        ("/example/foo/bar/baz", None)
    """

    paths = []
    while True:
        (mountpoint,), err = zfs.get(dset, ZFS.Property.MOUNTPOINT)
        if err != None:
            parent, base = os.path.split(dset)
            if parent == "":
                return "", err

            paths.append(base)
            dset = parent

            continue

        paths.reverse()
        return os.path.join(mountpoint, *paths), None


def run(module: AnsibleModule) -> tuple[dict, Error | None]:
    props, err = ZFS.parse_create_properties(module.params["props"])
    if err != None:
        return {}, err

    parent_props, err = ZFS.parse_create_properties(
        module.params["parent_props"]
    )
    if err != None:
        return {}, err

    name = module.params["name"]

    zfs = ZFS(module)

    (mountpoint,), err = zfs.get(name, ZFS.Property.MOUNTPOINT)
    if err != None:
        err = zfs.create(name, props, parent_props)
        if err != None:
            return {}, err

        (mountpoint,), err = (
            get_mountpoint_in_check_mode(zfs, name)
            if module.check_mode
            else zfs.get(name, ZFS.Property.MOUNTPOINT)
        )
        if err != None:
            return {}, err

    return (
        dict(
            name=name, mountpoint=mountpoint, changed=zfs.changed, diff=zfs.diff
        ),
        None,
    )


def main() -> None:
    spec_on_off_noauto = dict(type=str, choices=["on", "off", "noauto"])
    spec_on_off = dict(type=str, choices=["on", "off"])

    module = AnsibleModule(
        argument_spec=dict(
            name=dict(type="str", required=True),
            props=dict(
                type="dict",
                default={},
                options=dict(
                    canmount=spec_on_off_noauto,
                    readonly=spec_on_off,
                ),
            ),
            parent_props=dict(
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

    try:
        res, err = run(module)
        if err != None:
            module.fail_json(
                msg=str(err),
            )

        module.exit_json(**res)
    except Exception:
        module.fail_json(msg=traceback.format_exc())


if "__main__" == __name__:
    main()
