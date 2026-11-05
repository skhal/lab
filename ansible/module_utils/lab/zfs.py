# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
from enum import StrEnum, auto

from ansible.module_utils.basic import AnsibleModule

from ansible.module_utils.lab.errors import Error


class InvalidPropertyError(Error):
    """InvalidPropertyError means used property is invalid or unsupported."""

    def __init__(self, name: str) -> None:
        super().__init__(f"invalid property '{name}'")


class DatasetNotExistError(Error):
    """DatasetNotExistError means specified dataset does not exit."""

    def __init__(self, name: str) -> None:
        super().__init__(f"can't open '{name}': dataset does not exist")


class ZFS:
    class Property(StrEnum):
        # keep-sorted start
        CAN_MOUNT = "canmount"
        MOUNTPOINT = auto()
        NAME = auto()
        READ_ONLY = "readonly"
        # keep-sorted end

    class Command(StrEnum):
        # keep-sorted start
        CREATE = auto()
        GET = auto()
        # keep-sorted end

    class GetFlag(StrEnum):
        # keep-sorted start
        COLUMN = "-o"
        NO_HEADER = "-H"
        PARSABLE_VALUE = "-p"
        # keep-sorted end

    class CreateFlag(StrEnum):
        OPTION = "-o"

    _CMD = "/sbin/zfs"
    _ALLOWED_CREATE_PROPS = [Property.CAN_MOUNT, Property.READ_ONLY]

    def __init__(self, module: AnsibleModule) -> None:
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
        cmd = [
            ZFS._CMD,
            ZFS.Command.GET,
            ZFS.GetFlag.NO_HEADER,
            ZFS.GetFlag.PARSABLE_VALUE,
            ZFS.GetFlag.COLUMN,
            "value",
            prop,
            dset,
        ]

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

        if not self._module.check_mode:
            if (err := self._create(dset, props)) is not None:
                return err

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

    def _maybe_create_parent(
        self, dset: str, props: dict[Property, str]
    ) -> Error | None:
        head, _ = os.path.split(dset)

        _, err = self.get(ZFS.Property.NAME, head)
        if err != None:
            return self.create(head, props, props)

        return None

    def _create(self, dset: str, props: dict[Property, str]) -> Error | None:
        cmd = [ZFS._CMD, ZFS.Command.CREATE]
        for prop, val in props.items():
            cmd.extend((ZFS.CreateFlag.OPTION, f"{prop}={val}"))
        cmd.append(dset)

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(stderr)

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
