# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
from enum import StrEnum, auto
from typing import List

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
        ORIGIN = auto()
        READ_ONLY = "readonly"
        # keep-sorted end

    class Command(StrEnum):
        # keep-sorted start
        CLONE = auto()
        CREATE = auto()
        GET = auto()
        LIST = auto()
        # keep-sorted end

    class GetFlag(StrEnum):
        # keep-sorted start
        COLUMN = "-o"
        NO_HEADER = "-H"
        PARSABLE_VALUE = "-p"
        # keep-sorted end

    class CreateFlag(StrEnum):
        OPTION = "-o"

    class CloneFlag(StrEnum):
        OPTION = "-o"

    class ListFlag(StrEnum):
        # keep-sorted start
        NO_HEADER = "-H"
        PROPERTY = "-o"
        TYPE = "-t"
        # keep-sorted end

    class ListProperty(StrEnum):
        NAME = auto()

    class FSType(StrEnum):
        SNAPSHOT = auto()

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

    def get(
        self, dset: str, *props: Property
    ) -> tuple[List[str], Error | None]:
        cmd = [
            ZFS._CMD,
            ZFS.Command.GET,
            ZFS.GetFlag.NO_HEADER,
            ZFS.GetFlag.PARSABLE_VALUE,
            ZFS.GetFlag.COLUMN,
            "value",
            ",".join(props),
            dset,
        ]

        rc, stdout, _ = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return [], DatasetNotExistError(dset)

        return stdout.splitlines(), None

    def clone(
        self, snapshot: str, dset: str, props: dict[Property, str]
    ) -> Error | None:
        if not self._module.check_mode:
            if (err := self._clone(snapshot, dset, props)) is not None:
                return err

        self._changed = True
        self._diff.append(
            {
                "before_header": dset,
                "before": dict(state="absent"),
                "after_header": dset,
                "after": dict(state="present", origin=snapshot),
            }
        )

    def _clone(
        self, snapshot: str, dset: str, props: dict[Property, str]
    ) -> Error | None:
        cmd = [ZFS._CMD, ZFS.Command.CLONE]
        for prop, val in props.items():
            cmd.extend((ZFS.CloneFlag.OPTION, f"{prop}={val}"))
        cmd.extend((snapshot, dset))

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            if "dataset does not exist" in stderr:
                return DatasetNotExistError(snapshot)
            return Error(stderr.strip())

        return None

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

        _, err = self.get(head, ZFS.Property.NAME)
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

    def list(self, dset: str, type: FSType) -> tuple[List[str], Error | None]:
        cmd = [
            ZFS._CMD,
            ZFS.Command.LIST,
            # keep-sorted start newline_separated=yes group_start_regex=^\s+ZFS\.ListFlag\.
            ZFS.ListFlag.NO_HEADER,

            ZFS.ListFlag.PROPERTY,
            ZFS.ListProperty.NAME,

            ZFS.ListFlag.TYPE,
            type,
            # keep-sorted end
            dset,
        ]

        rc, stdout, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return [], Error(stderr.strip())

        return stdout.splitlines(), None

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
