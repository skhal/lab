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
    rootdir: /tmp
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


class PW:
    """Runs pw(8) to rootdir."""

    # Ref: https://github.com/freebsd/freebsd-src/blob/7f5fa76367d78e47d483fdf2cc72e5823d0f7807/lib/libutil/pw_util.c#L400-L403
    _SHOWUSER_IDX_NAME = 0
    _SHOWUSER_IDX_UID = 2
    _SHOWUSER_IDX_GID = 3

    def __init__(self, module: AnsibleModule, rootdir: str):
        self._module = module
        self._rootdir = rootdir

    def usershow(self, name: str) -> User:
        """Runs pw(8) usershow command to extract user information."""
        rc, out, err = self._module.run_command(
            args=["pw", "-R", self._rootdir, "usershow", name],
        )
        if rc == os.EX_NOUSER:
            raise NoUserError()
        if rc != 0:
            raise OSError(rc, err)
        tokens = out.split(":")
        return User(
            name=tokens[self._SHOWUSER_IDX_NAME],
            uid=tokens[self._SHOWUSER_IDX_UID],
            gid=tokens[self._SHOWUSER_IDX_GID],
        )


class UserMigrator:
    """Migrate user to the rootdir"""

    def __init__(self, module):
        self._module = module

    def migrate(self, name: str, rootdir: str) -> dict:
        """Migrates user with a given name to the rootdir."""
        try:
            PW(self._module, rootdir).usershow(name)
            return dict(
                msg=f"user {name} exists",
                changed=False,
            )
        except NoUserError:
            pass
        except Exception as err:
            return dict(
                msg=f"failed to migrate user {name}\n{err}",
                failed=True,
            )
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
    try:
        run(module)
    except Exception as ex:
        module.fail_json(msg="failed to migrate users", exception=ex)


if __name__ == "__main__":
    main()
