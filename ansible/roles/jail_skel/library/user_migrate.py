# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

import os
import traceback
from dataclasses import dataclass
from enum import IntEnum, StrEnum

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.lab.errors import Error


class UserNotFoundError(Error):
    def __init__(self, *users: str) -> None:
        super().__init__(f"users '{','.join(users)}' not found")


@dataclass(eq=True, frozen=True)
class User:
    name: str
    uid: int
    gid: int
    gecos: str
    home: str
    shell: str


# passwd record:
#   % pw usershow root
#   root:*:0:0::0:0:Charlie &:/root:/bin/tcsh
#
# ref: https://github.com/freebsd/freebsd-src/blob/9145e8a2ebe57b260e8e929043ed2052188eddae/lib/libutil/pw_util.c#L400-L403
class PasswdIndex(IntEnum):
    NAME = 0
    UID = 2
    GID = 3
    GECOS = 7
    HOME = 8
    SHELL = 9


def parse_passwd_record(record: str) -> User:
    tokens = record.split(sep=":")
    return User(
        name=tokens[PasswdIndex.NAME],
        uid=int(tokens[PasswdIndex.UID]),
        gid=int(tokens[PasswdIndex.GID]),
        gecos=tokens[PasswdIndex.GECOS],
        home=tokens[PasswdIndex.HOME],
        shell=tokens[PasswdIndex.SHELL],
    )


@dataclass(eq=True, frozen=True)
class Group:
    name: str
    gid: int
    users: frozenset[str]


# group record:
#   % pw groupshow wheel
#   wheel:*:0:root,user
#
# ref: https://github.com/freebsd/freebsd-src/blob/208d30a8f594e5683992a71dd01747179a1f923d/lib/libutil/gr_util.c#L448-L457
class GroupIndex(IntEnum):
    NAME = 0
    GID = 2
    USERS = 3


def parse_group_record(record: str) -> Group:
    tokens = record.split(sep=":")
    group = Group(
        name=tokens[GroupIndex.NAME],
        gid=int(tokens[GroupIndex.GID]),
        users=frozenset(tokens[GroupIndex.USERS].split(",")),
    )
    return group


class PW:
    _CMD = "/usr/sbin/pw"

    class Command(StrEnum):
        # keep-sorted start by_regex="(\w+)"
        GroupAdd = "groupadd"
        GroupModify = "groupmod"
        GroupShow = "groupshow"
        UserAdd = "useradd"
        UserModify = "usermod"
        UserShow = "usershow"
        # keep-sorted end

    class Flag(StrEnum):
        Rootdir = "-R"

    # UserFlags lists used flags from pw(8) for useradd.
    class UserFlag(StrEnum):
        # keep-sorted start by_regex=(-\w)
        GroupList = "-G"
        Comment = "-c"
        HomeDir = "-d"
        GroupID = "-g"
        CreateUserHome = "-m"
        Name = "-n"
        Shell = "-s"
        UserID = "-u"
        # keep-sorted end

    class GroupFlag(StrEnum):
        # keep-sorted start by_regex="(-\w)"
        NewMembers = "-M"
        GroupID = "-g"
        Name = "-n"
        # keep-sorted end

    def __init__(
        self, module: AnsibleModule, *, rootdir: str | None = None
    ) -> None:
        self._module = module
        self._rootdir = rootdir
        self._changed = False
        self._diff = []

    @property
    def changed(self):
        return self._changed

    @property
    def diff(self):
        return self._diff

    def useradd(self, u: User, supplemental_groups: list[int]) -> Error | None:
        if not self._module.check_mode:
            if (err := self._useradd(u, supplemental_groups)) is not None:
                return err

        self._changed = True
        self._diff.append(
            {
                "before_header": f"user: {u.name}",
                "before": dict(state="absent"),
                "after_header": f"user: {u.name}",
                "after": dict(state="present"),
            }
        )

    def _useradd(self, u: User, supplemental_groups: list[int]) -> Error | None:
        cmd = [PW._CMD]
        if self._rootdir is not None:
            cmd.extend([PW.Flag.Rootdir, self._rootdir])
        cmd.extend(
            [
                PW.Command.UserAdd,
                # keep-sorted start group_start_regex=['^\s+PW\.UserFlags']
                PW.UserFlag.Comment,
                u.gecos,
                PW.UserFlag.CreateUserHome,
                PW.UserFlag.GroupID,
                str(u.gid),
                PW.UserFlag.HomeDir,
                u.home,
                PW.UserFlag.Name,
                u.name,
                PW.UserFlag.Shell,
                u.shell,
                PW.UserFlag.UserID,
                str(u.uid),
                # keep-sorted end
            ]
        )
        if supplemental_groups:
            cmd.extend(
                [
                    PW.UserFlag.GroupList,
                    ",".join(str(gid) for gid in supplemental_groups),
                ]
            )

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(
                f"{PW.Command.UserAdd}: {stderr.strip()}\n{u}"
                f" groups:{supplemental_groups}"
            )

    def usermod(self, u: User) -> Error | None:
        if not self._module.check_mode:
            if (err := self._usermod(u)) is not None:
                return err

        self._changed = True
        self._diff.append(
            {
                "before_header": f"user: {u.name}",
                "before": dict(state="present"),
                "after_header": f"user: {u.name}",
                "after": dict(state="updated"),
            }
        )

    def _usermod(self, u: User) -> Error | None:
        cmd = [PW._CMD]
        if self._rootdir is not None:
            cmd.extend([PW.Flag.Rootdir, self._rootdir])
        cmd.extend(
            [
                PW.Command.UserModify,
                # keep-sorted start group_start_regex=['^\s+PW\.UserFlags']
                PW.UserFlag.Comment,
                u.gecos,
                PW.UserFlag.CreateUserHome,
                PW.UserFlag.GroupID,
                str(u.gid),
                PW.UserFlag.HomeDir,
                u.home,
                PW.UserFlag.Name,
                u.name,
                PW.UserFlag.Shell,
                u.shell,
                PW.UserFlag.UserID,
                str(u.uid),
                # keep-sorted end
            ]
        )

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(f"{PW.Command.UserModify}: {stderr.strip()}\n{u}")

    def usershow_all(self) -> tuple[list[User], Error | None]:
        cmd = [PW._CMD]
        if self._rootdir is not None:
            cmd.extend([PW.Flag.Rootdir, self._rootdir])
        cmd.extend([PW.Command.UserShow, "-a"])

        rc, stdout, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return [], Error(f"{PW.Command.UserShow}: {stderr.strip()}")

        users = list(parse_passwd_record(line) for line in stdout.splitlines())

        return users, None

    def groupadd(self, g: Group) -> Error | None:
        if not self._module.check_mode:
            if (err := self._groupadd(g)) is not None:
                return err

        self._changed = True
        self._diff.append(
            {
                "before_header": f"group: {g.name}",
                "before": dict(state="absent"),
                "after_header": f"group: {g.name}",
                "after": dict(state="present"),
            }
        )

    def _groupadd(self, g: Group) -> Error | None:
        cmd = [PW._CMD]
        if self._rootdir is not None:
            cmd.extend([PW.Flag.Rootdir, self._rootdir])
        cmd.extend(
            [
                PW.Command.GroupAdd,
                # keep-sorted start group_start_regex=['^\s+PW\.GroupFlags']
                PW.GroupFlag.GroupID,
                str(g.gid),
                PW.GroupFlag.Name,
                g.name,
                # keep-sorted end
            ]
        )

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(f"{PW.Command.GroupAdd}: {stderr.strip()}\n{g}")

    def groupmod(self, g: Group) -> Error | None:
        if not self._module.check_mode:
            if (err := self._groupmod(g)) is not None:
                return err

        self._changed = True
        self._diff.append(
            {
                "before_header": f"group: {g.name}",
                "before": dict(state="present"),
                "after_header": f"group: {g.name}",
                "after": dict(state="updated"),
            }
        )

    def _groupmod(self, g: Group) -> Error | None:
        cmd = [PW._CMD]
        if self._rootdir is not None:
            cmd.extend([PW.Flag.Rootdir, self._rootdir])
        cmd.extend(
            [
                PW.Command.GroupModify,
                # keep-sorted start group_start_regex=['^\s+PW\.GroupFlags']
                PW.GroupFlag.GroupID,
                str(g.gid),
                PW.GroupFlag.Name,
                g.name,
                # want to preserve existing group members
                PW.GroupFlag.NewMembers,
                ",".join(g.users),
                # keep-sorted end
            ]
        )

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(f"{PW.Command.GroupModify}: {stderr.strip()}\n{g}")

    def groupshow_all(self) -> tuple[list[Group], Error | None]:
        cmd = [PW._CMD]
        if self._rootdir is not None:
            cmd.extend([PW.Flag.Rootdir, self._rootdir])
        cmd.extend([PW.Command.GroupShow, "-a"])

        rc, stdout, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return [], Error(f"{PW.Command.GroupShow}: {stderr.strip()}")

        groups = list(parse_group_record(line) for line in stdout.splitlines())

        return groups, None


@dataclass(frozen=True)
class AdminDB:
    users: list[User]
    groups: list[Group]

    def empty(self) -> bool:
        return not self.users and not self.groups


class UserMigrator:
    def __init__(self, module: AnsibleModule) -> None:
        self._module = module
        self._migrated_users = []
        self._migrated_groups = []
        self._changed = False
        self._diff = []

    @property
    def migrated_users(self) -> list[str]:
        return self._migrated_users

    @property
    def migrated_groups(self) -> list[str]:
        return self._migrated_groups

    @property
    def changed(self) -> bool:
        return self._changed

    @property
    def diff(self) -> list[dict]:
        return self._diff

    def migrate(self, users: list[str], dest: str) -> Error | None:
        if not users:
            return

        src_admdb, err = self._getent()
        if err is not None:
            return err

        if src_admdb.empty():
            return Error(f"users '{','.join(users)}' not found")

        # all requested users must be present at the source.
        src_admdb, src_missing_users = self._filter(src_admdb, users)
        if src_missing_users:
            return UserNotFoundError(*src_missing_users)

        dst_admdb = AdminDB([], [])
        if os.path.exists(dest):
            dst_admdb, err = self._getent(dest)
            if err is not None:
                return err

        add_admdb, mod_admdb = self._difference(src_admdb, dst_admdb)

        return self._migrate(add_admdb, mod_admdb, dest)

    def _getent(
        self, rootdir: str | None = None
    ) -> tuple[AdminDB, Error | None]:
        pw = PW(self._module, rootdir=rootdir)

        users, err = pw.usershow_all()
        if err is not None:
            return AdminDB([], []), err

        groups, err = pw.groupshow_all()
        if err is not None:
            return AdminDB([], []), err

        return AdminDB(users=users, groups=groups), None

    def _filter(
        self, db: AdminDB, users: list[str]
    ) -> tuple[AdminDB, list[str]]:
        if db.empty():
            return db, []

        requested_users = set(users)

        filtered_users = list(u for u in db.users if u.name in requested_users)
        filtered_user_groups = set(u.gid for u in filtered_users)

        missing_users = requested_users - set(u.name for u in filtered_users)

        def is_keep_group(g: Group) -> bool:
            if g.users & requested_users:
                return True

            return g.gid in filtered_user_groups

        filtered_groups = list(g for g in db.groups if is_keep_group(g))

        return AdminDB(filtered_users, filtered_groups), sorted(missing_users)

    def _difference(self, x: AdminDB, y: AdminDB) -> tuple[AdminDB, AdminDB]:
        add_users, mod_users = self._gen_diff(x.users, y.users)
        add_groups, mod_groups = self._gen_diff(x.groups, y.groups)

        add_db = AdminDB(add_users, add_groups)
        mod_db = AdminDB(mod_users, mod_groups)

        return add_db, mod_db

    def _gen_diff[T: User | Group](
        self, x: list[T], y: list[T]
    ) -> tuple[list[T], list[T]]:
        x_set = set(x)
        y_set = set(y)

        common = x_set & y_set

        x_set -= common
        y_set -= common

        y_names_set = set(v.name for v in y_set)

        add_from_x = list(v for v in x_set if v.name not in y_names_set)
        mod_from_x = list(v for v in x_set if v.name in y_names_set)

        return add_from_x, mod_from_x

    def _migrate(
        self, add_db: AdminDB, mod_db: AdminDB, rootdir: str
    ) -> Error | None:
        if (err := self._add(add_db, rootdir)) is not None:
            return err

        return self._modify(mod_db, rootdir)

    def _add(self, db: AdminDB, rootdir: str) -> Error | None:
        if db.empty():
            return

        pw = PW(self._module, rootdir=rootdir)

        group_users = {}
        for g in db.groups:
            for u in g.users:
                group_users.setdefault(u, set()).add(g.gid)
            if (err := pw.groupadd(g)) is not None:
                return err

        self._migrated_groups.extend(g.name for g in db.groups)

        for u in db.users:
            supplemental_groups = group_users.get(u.name, set())
            if u.gid in supplemental_groups:
                supplemental_groups.remove(u.gid)

            if (err := pw.useradd(u, supplemental_groups)) is not None:
                return err

        self._migrated_users.extend(u.name for u in db.users)

        self._changed |= pw.changed
        self._diff.extend(pw.diff)

    def _modify(self, db: AdminDB, rootdir: str) -> Error | None:
        if db.empty():
            return

        pw = PW(self._module, rootdir=rootdir)

        for g in db.groups:
            if (err := pw.groupmod(g)) is not None:
                return err

        self._migrated_groups.extend(g.name for g in db.groups)

        for u in db.users:
            if (err := pw.usermod(u)) is not None:
                return err

        self._migrated_users.extend(u.name for u in db.users)

        self._changed |= pw.changed
        self._diff.extend(pw.diff)


def run(module: AnsibleModule) -> tuple[dict, Error | None]:
    users = module.params["name"]
    dest = module.params["dest"]

    um = UserMigrator(module)

    if (err := um.migrate(users, dest)) is not None:
        return {}, err

    return (
        dict(
            users=um.migrated_users,
            groups=um.migrated_groups,
            changed=um.changed,
            diff=um.diff,
        ),
        None,
    )


def main() -> None:
    module = AnsibleModule(
        argument_spec=dict(
            name=dict(type="list", elements="str", required=True),
            dest=dict(type="str", required=True),
        ),
        supports_check_mode=True,
    )

    try:
        res, err = run(module)
        if err is not None:
            module.fail_json(msg=str(err))

        module.exit_json(**res)
    except Exception:
        module.fail_json(msg=traceback.format_exc())


if "__main__" == __name__:
    main()
