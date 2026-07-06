# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

DOCUMENTATION = r"""
---
module: migrate
short_description: Migrate users from / to rootdir

options:
    name:
        description: User name to migrate.
        type: str
        required: true
    rootdir:
        description: Destination rootdir to migrate the user to.
        type: str
        required: true

author: Samvel Khalatyan (sn.khalatyan@gmail.com)
"""

EXAMPLES = r"""
# Migrate user op from / to the new rootdir /usr/local/foo
- name: Migrate user op to /usr/local/foo
  migrate:
    name: op
    rootdir: /usr/local/foo
"""

import os
from dataclasses import dataclass, field
from typing import NewType, Tuple
from ansible.module_utils.basic import AnsibleModule


class NoUserError(Exception):
    """Corresponds to os.EX_NOUSER error."""


class GroupError(Exception):
    """Means group mismatch between the source and target rootdirs."""


@dataclass
class User:
    """Describes a user in passwd database."""

    name: str
    uid: int
    gid: int
    gecos: str
    home: str
    shell: str
    secondary_gids: list[int] = field(default_factory=list)

    def with_groups(self, gids: list[int]):
        """Duplicates current user with secondary groups include gids.

        It ensures that user primary group is not included in the secondary
        groups.
        """
        groups = set(gids).union(self.secondary_gids) - {self.gid}
        return User(
            name=self.name,
            uid=self.uid,
            gid=self.gid,
            gecos=self.gecos,
            home=self.home,
            shell=self.shell,
            secondary_gids=list(groups),
        )

    def groups(self) -> list[int]:
        return [self.gid] + self.secondary_gids


@dataclass
class Group:
    """Describes a group in groups databasae."""

    name: str
    gid: int


class PW:
    """Runs pw(8) in rootdir."""

    _DEFAULT_ROOTDIR = "/"

    # Ref: https://github.com/freebsd/freebsd-src/blob/7f5fa76367d78e47d483fdf2cc72e5823d0f7807/lib/libutil/pw_util.c#L400-L403
    _USERSHOW_IDX_NAME = 0
    _USERSHOW_IDX_UID = 2
    _USERSHOW_IDX_GID = 3
    _USERSHOW_IDX_GECOS = 7
    _USERSHOW_IDX_HOME = 8
    _USERSHOW_IDX_SHELL = 9

    _GROUPSHOW_IDX_NAME = 0
    _GROUPSHOW_IDX_GID = 2

    def __init__(self, module: AnsibleModule, rootdir: str = ""):
        self._module = module
        self._rootdir = rootdir if rootdir else self._DEFAULT_ROOTDIR

    def groupadd(self, group: Group):
        rc, _, err = self._module.run_command(
            args=[
                "pw",
                "-R",
                self._rootdir,
                "groupadd",
                "-n",
                group.name,
                "-g",
                str(group.gid),
            ],
        )
        if rc != os.EX_OK:
            raise OSError(rc, err)

    def groupshow_all(self) -> dict[int, Group]:
        rc, out, err = self._module.run_command(
            args=["pw", "-R", self._rootdir, "groupshow", "-a"],
        )
        match rc:
            case os.EX_OK:
                pass
            case _:
                raise OSError(rc, err)
        return self._parse_groupshow(out)

    def _parse_groupshow(self, out: str) -> dict[int, Group]:
        groups = dict()
        for line in out.splitlines():
            tokens = line.split(":")
            group = Group(
                name=tokens[self._GROUPSHOW_IDX_NAME],
                gid=int(tokens[self._GROUPSHOW_IDX_GID]),
            )
            groups[group.gid] = group
        return groups

    def useradd(self, user: User):
        rc, _, err = self._module.run_command(
            args=[
                "pw",
                "-R",
                self._rootdir,
                "useradd",
                "-n",
                user.name,
                "-u",
                str(user.uid),
                "-g",
                str(user.gid),
                "-c",
                user.gecos,
                "-d",
                user.home,
                "-s",
                user.shell,
            ],
        )
        if rc != os.EX_OK:
            raise OSError(rc, err)

    def usershow(self, name: str) -> User:
        """Runs pw(8) usershow command to extract user information."""
        rc, out, err = self._module.run_command(
            args=["pw", "-R", self._rootdir, "usershow", name],
        )
        match rc:
            case os.EX_NOUSER:
                raise NoUserError()
            case os.EX_OK:
                pass
            case _:
                raise OSError(rc, err)
        tokens = out.split(":")
        return User(
            name=tokens[self._USERSHOW_IDX_NAME],
            uid=int(tokens[self._USERSHOW_IDX_UID]),
            gid=int(tokens[self._USERSHOW_IDX_GID]),
            gecos=tokens[self._USERSHOW_IDX_GECOS],
            home=tokens[self._USERSHOW_IDX_HOME],
            shell=tokens[self._USERSHOW_IDX_SHELL],
        )


class UserIdentity:
    """Wraps id(1)."""

    def __init__(self, module: AnsibleModule):
        self._module = module

    def groups(self, name: str) -> list[int]:
        """Gets a list of user groups including primary and secondary groups."""
        rc, out, err = self._module.run_command(
            args=["id", "-G", name],
        )
        if rc != os.EX_OK:
            raise OSError(rc, err)
        return list(int(n) for n in out.split())


class UserMigrator:
    """Migrate user to the rootdir"""

    def __init__(self, module):
        self._module = module

    def migrate(self, name: str, rootdir: str) -> dict:
        """Migrates user with a given name to the rootdir."""
        try:
            return self._migrate(name, rootdir)
        except Exception as ex:
            return dict(
                msg=f"failed to migrate user {name}\n{ex}",
                failed=True,
            )

    def _migrate(self, name: str, rootdir: str) -> dict:
        try:
            PW(self._module, rootdir).usershow(name)
            return dict(
                msg=f"user {name} exists",
                changed=False,
            )
        except NoUserError:
            pass
        user = self._get_user(name)
        self._migrate_groups(user, rootdir)
        self._migrate_user(user, rootdir)
        return dict(msg=f"user {name} migrated", changed=True)

    def _get_user(self, name: str) -> User:
        user = PW(self._module).usershow(name)
        return user.with_groups(UserIdentity(self._module).groups(name))

    def _migrate_groups(self, user: User, rootdir: str):
        src_groups = PW(self._module).groupshow_all()
        dst_groups = PW(self._module, rootdir).groupshow_all()

        for gid in user.groups():
            group = src_groups[gid]
            if self._validate_group(group, dst_groups.get(gid, None)):
                continue
            self._create_group(group, rootdir)

    def _validate_group(self, src: Group, dst: Group | None):
        if dst == None:
            return False
        if src != dst:
            raise GroupError(
                f"group {src.gid} name {dst.name} does not match {src.name}"
            )
        return True

    def _create_group(self, group: Group, rootdir: str):
        PW(self._module, rootdir).groupadd(group)

    def _migrate_user(self, user: User, rootdir: str):
        PW(self._module, rootdir).useradd(user)


def run(module: AnsibleModule):
    if module.check_mode:
        module.exit_json(changed=False)
    result = UserMigrator(module).migrate(
        name=module.params["name"],
        rootdir=module.params["rootdir"],
    )
    if "failed" in result:
        module.fail_json(**result)
    else:
        module.exit_json(**result)


def main():
    module = AnsibleModule(
        argument_spec=dict(
            name=dict(
                type="str",
                required=True,
            ),
            rootdir=dict(
                type="str",
                required=True,
            ),
        ),
        supports_check_mode=True,
    )
    run(module)


if __name__ == "__main__":
    main()
