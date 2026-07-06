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
from dataclasses import dataclass
from typing import NewType, Tuple
from ansible.module_utils.basic import AnsibleModule


class NoUserError(Exception):
    """Corresponds to os.EX_NOUSER error."""

    pass


@dataclass
class User:
    """Describes a user in passwd database."""

    name: str
    uid: int
    gid: int
    gecos: str
    home: str
    shell: str


class PW:
    """Runs pw(8) to rootdir."""

    _DEFAULT_ROOTDIR = "/"

    # Ref: https://github.com/freebsd/freebsd-src/blob/7f5fa76367d78e47d483fdf2cc72e5823d0f7807/lib/libutil/pw_util.c#L400-L403
    _SHOWUSER_IDX_NAME = 0
    _SHOWUSER_IDX_UID = 2
    _SHOWUSER_IDX_GID = 3
    _SHOWUSER_IDX_GECOS = 7
    _SHOWUSER_IDX_HOME = 8
    _SHOWUSER_IDX_SHELL = 9

    def __init__(self, module: AnsibleModule, rootdir: str = ""):
        self._module = module
        self._rootdir = rootdir if rootdir else self._DEFAULT_ROOTDIR

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
                user.uid,
                "-g",
                user.gid,
                # TODO(github.com/skhal/lab/issues/400): add supplemental groups
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
            name=tokens[self._SHOWUSER_IDX_NAME],
            uid=tokens[self._SHOWUSER_IDX_UID],
            gid=tokens[self._SHOWUSER_IDX_GID],
            gecos=tokens[self._SHOWUSER_IDX_GECOS],
            home=tokens[self._SHOWUSER_IDX_HOME],
            shell=tokens[self._SHOWUSER_IDX_SHELL],
        )


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
        user = PW(self._module).usershow(name)
        # TODO(github.com/skhal/lab/issues/400): migrate groups
        PW(self._module, rootdir).useradd(user)
        return dict(msg=f"user {name} migrated", changed=False)


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
