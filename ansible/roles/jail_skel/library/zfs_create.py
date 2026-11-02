# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

DOCUMENTATION = r"""
module: zfs_create
short_description: Create a ZFS dataset and its parents
description:
    - Create a ZFS dataset along with parents with different properties.
attributes:
    check_mode:
        support: full
    diff_mode:
        support: full
options:
    name:
        description:
            - ZFS dataset to create
        type: str
        required: true
    props:
        description:
            - ZFS dataset properties
        type: dict
        default: {}
    parent_props:
        description:
            - parent dataset properties
        type: dict
        default: {}
author: Samvel Khalatyan (@skhal)
"""

EXAMPLES = r"""
# Setup: zpool/example dataset exists

# Create a dataset and parents with default properties
- name: Create a dataset and parents
  zfs_create:
    name: zpool/example/dot/com

# Create a leaf dataset with custom properties:
- name: Create a dataset and parents
  zfs_create:
    name: zpool/example/dot/com
    props:
        readonly: "no"

# Create a leaf dataset and parent datasets with different properties
- name: Create a dataset and parents
  zfs_create:
    name: zpool/example/dot/com
    props:
        readonly: "no"
    parent_props:
        readonly: "yes"
"""

RETURN = r"""
name:
    description: Dataset name
    returned: always
    type: str
    sample: zroot/usr/example
mountpoint:
    description: Dataset mountpoint
    returned: always
    type: str
    sample: /usr/example
"""

import os

from ansible.module_utils.basic import AnsibleModule


class Error:
    def __init__(self, msg: str):
        self._msg = msg

    def __str__(self):
        return self._msg


class DatasetNotExistError(Error):
    """DatasetNotExistError means specified dataset does not exit."""

    def __init__(self, name: str):
        super().__init__(f"can't open '{name}': dataset does not exist")


class InvalidPropertyError(Error):
    """InvalidPropertyError means used property is invalid or unsupported."""

    def __init__(self, name: str):
        super().__init__(f"invalid property '{name}'")


class Property:
    """Property describes a ZFS dataset property.

    Derived classes are supported properties. Use these instead of strings to
    validate supported properties and use it in annotations.

    This implementation only supports a handful of zfsprops(7) that are used
    in current module.
    """

    def __init__(self, name: str):
        self._name = name

    def __str__(self):
        return self._name


class CanMountProperty(Property):
    def __init__(self):
        super().__init__("canmount")


class MountpointProperty(Property):
    def __init__(self):
        super().__init__("mountpoint")


class NameProperty(Property):
    def __init__(self):
        super().__init__("name")


class ReadOnlyProperty(Property):
    def __init__(self):
        super().__init__("readonly")


class ZFS:
    _ZFS_CMD = "/sbin/zfs"
    _ALLOWED_CREATE_PROPS = [CanMountProperty(), ReadOnlyProperty()]

    def __init__(self, module: AnsibleModule):
        self._module = module
        self._changed = False
        self._diff = []

    @property
    def changed(self):
        return self._changed

    @property
    def diff(self):
        return self._diff

    def get(self, prop: Property, dset: str) -> tuple[str, Error | None]:
        cmd = [ZFS._ZFS_CMD, "get", "-H", "-p", "-o", "value", str(prop), dset]

        rc, stdout, _ = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return "", DatasetNotExistError(dset)

        return stdout.strip(), None

    def create(
        self,
        dset: str,
        props: dict[Property, str] = {},
        parent_props: dict[Property, str] = {},
    ) -> Error | None:
        """create implements zfs-create(8).

        It creates a dataset with optional properties. It also creates missing
        parent datasets with different optional properties.

        The difference in leaf and parent dataset properties unlocks parent
        datasets that are commonly used to propagate properties such as
        mountpoints but should not be mountable themselves.
        """

        err = self._maybe_create_parent(dset, parent_props)
        if err != None:
            return err

        return self._create(dset, props)

    def _maybe_create_parent(
        self, dset: str, props: dict[Property, str]
    ) -> Error | None:
        head, _ = os.path.split(dset)

        _, err = self.get(NameProperty(), head)
        if err != None:
            return self.create(head, props, props)

        return None

    def _create(self, dset: str, props: dict[Property, str]) -> Error | None:

        if not self._module.check_mode:
            cmd = [ZFS._ZFS_CMD, "create"]
            for prop, val in props.items():
                cmd.extend(("-o", f"{prop}={val}"))
            cmd.append(dset)

            rc, _, stderr = self._module.run_command(cmd)
            if rc != os.EX_OK:
                return Error(stderr)

        self._changed = True
        self._diff.append(
            {
                "before_header": dset,
                "before": dict(state="absent"),
                "after_header": dset,
                # convert props to a dict[str,str] for JSON
                "after": dict(
                    state="present", props={str(p): v for p, v in props.items()}
                ),
            }
        )

        return None

    @staticmethod
    def parse_create_properties(
        props: dict,
    ) -> tuple[dict[Property, str], Error | None]:
        """converts ZFS properties from strings to strong types."""

        propByName = dict({str(p): p for p in ZFS._ALLOWED_CREATE_PROPS})
        parsed_props = {}
        for name, val in props.items():
            prop = propByName.get(name)
            if prop is None:
                return {}, InvalidPropertyError(name)
            parsed_props[prop] = val
        return parsed_props, None


def get_mountpoint_in_check_mode(zfs: ZFS, dset: str) -> tuple[str, Error | None]:
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
        mountpoint, err = zfs.get(MountpointProperty(), dset)
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

    parent_props, err = ZFS.parse_create_properties(module.params["parent_props"])
    if err != None:
        return {}, err

    name = module.params["name"]

    zfs = ZFS(module)

    mountpoint, err = zfs.get(MountpointProperty(), name)
    if err != None:
        err = zfs.create(name, props, parent_props)
        if err != None:
            return {}, err

        mountpoint, err = (
            get_mountpoint_in_check_mode(zfs, name)
            if module.check_mode
            else zfs.get(MountpointProperty(), name)
        )
        if err != None:
            return {}, err

    return (
        dict(name=name, mountpoint=mountpoint, changed=zfs.changed, diff=zfs.diff),
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

    res, err = run(module)
    if err != None:
        module.fail_json(
            msg=str(err),
        )
    else:
        module.exit_json(**res)


if "__main__" == __name__:
    main()
