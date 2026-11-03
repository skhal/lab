# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

DOCUMENTATION = r"""
module: jail_bootstrap
short_description: Install packages into empty, non-running jail
description:
    - Copy canonical pkg-repo(8) keys from the server into the jail.
    - Mount /var/db/pkg/repos as tmpfs in the jail.
    - Install packages from the server using `pkg --rootfs`
attributes:
    platform:
        platforms: FreeBSD
    check_mode:
        support: full
    diff_mode:
        support: full
options:
    path:
        description:
            - Jail location on the filesystem
        type: str
        required: true
    packages:
        description:
            - A list of package names to install
            - Do not include package version
        type: list
        required: true
requirements:
    - jail path must exist
    - /sbin/mount (pkg: FreeBSD-runtime)
    - /sbin/umount (pkg: FreeBSD-runtime)
    - /usr/sbin/pkg (pkg: FreeBSD-pkg-bootstrap)
author: Samvel Khalatyan (@skhal)
"""

EXAMPLES = r"""
# Setup: a jail path /usr/jail/example exists

# Install minimal jail set
- name: Install packages
  jail_bootstrap:
    path: /usr/jail/example
    packages:
        - FreeBSD-set-minimal-jail
        - doas
"""

RETURN = r"""
path:
    description: Jail location
    returned: always
    type: str
    sample: /usr/jail/example
packages:
    description: List of installed packaged
    returned: always
    type: list
    sample: [FreeBSD-set-minimal-jail, doas]
"""

import os
import shutil
import traceback
from contextlib import contextmanager

from ansible.module_utils.basic import AnsibleModule


class Error:
    def __init__(self, msg: str):
        self._msg = msg

    def __str__(self) -> str:
        return self._msg


class TMPFSContext:
    _MOUNT_CMD = "/sbin/mount"
    _UMOUNT_CMD = "/sbin/umount"

    def __init__(self, module: AnsibleModule, path: str):
        self._module = module
        self._path = path
        self._err = None

    @property
    def err(self) -> Error | None:
        return self._err

    def __enter__(self) -> TMPFSContext:
        self._err = self._mount()
        return self

    def _mount(self) -> Error | None:
        if self._module.check_mode:
            return

        cmd = [TMPFSContext._MOUNT_CMD, "-t", "tmpfs", "tmpfs", self._path]

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(stderr.strip())

    def __exit__(self, *_) -> bool:
        # Ansible modules are restricted to use of ansible.module_utils and
        # can't access ansible.errors.AnsibleError on the remote host unless
        # it is installed as external dependency and this way violate hermetic
        # setup of the module.
        #
        # In order to alleviate the above restriction, the module is written
        # exceptions free with explicit error propagation, e.g. Google C++ Style
        # Guide, Go, etc.
        #
        # Python libraries do use exceptions. The code should aggressively catch
        # these to convert exceptions into errors.
        #
        # Nevertheless some parts of the code may still miss exception handling.
        # Let it yell loud and clear with violent crash for easy detection.
        #
        # Propagate exception from the context block, catch it at the very high
        # level of module use, and include full traceback in the modules's error
        # message.

        if self._err != None:
            return False

        self._err = self._umount()

        return False

    def _umount(self) -> Error | None:
        if self._module.check_mode:
            return

        cmd = [TMPFSContext._UMOUNT_CMD, self._path]

        rc, _, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return Error(stderr.strip())


class PKG:
    _DB_REPOS_PATH = "/var/db/pkg/repos"
    _KEYS_PATH = "/usr/share/keys"

    _PKG_CMD = "/usr/sbin/pkg"

    def __init__(self, module: AnsibleModule):
        self._module = module
        self._installed_packages = []
        self._changed = False
        self._diff = []

    @property
    def installed_packages(self):
        return self._installed_packages

    @property
    def changed(self):
        return self._changed

    @property
    def diff(self):
        return self._diff

    def bootstrap(self, rootdir: str) -> Error | None:
        """bootstrap prepares rootdir to run pkg"""

        return self._install_keys(rootdir)

    def _install_keys(self, rootdir: str) -> Error | None:
        _, _, tail = os.path.splitroot(PKG._KEYS_PATH)
        dest = os.path.join(rootdir, tail)

        if not os.path.exists(dest):
            self._mkdir(dest)

        self._copytree(dest)

    def _mkdir(self, path: str) -> Error | None:
        if not self._module.check_mode:
            try:
                os.mkdir(path=path, mode=0o755)
            except OSError as err:
                return Error(str(err))

        self._changed = True
        self._diff.append(
            {
                "before_header": path,
                "before": dict(state="absent"),
                "after_header": path,
                "after": dict(state="present"),
            }
        )

    def _copytree(self, dest: str) -> Error | None:
        if self._module.check_mode:
            return

        copied_files = []

        def copy(src: str, dst: str, *, follow_symlinks=True):
            shutil.copy2(src, dst, follow_symlinks=follow_symlinks)
            copied_files.append(dst)

        try:
            shutil.copytree(PKG._KEYS_PATH, dest, copy_function=copy)
        except OSError as err:
            return Error(str(err))

        if not copied_files:
            return

        self._changed = True
        self._diff.extend(
            {
                "before_header": f,
                "before": dict(state="absent"),
                "after_header": f,
                "after": dict(state="present"),
            }
            for f in copied_files
        )

    def install(self, packages: list[str], *, rootdir: str) -> Error | None:
        """install listed packaged"""

        if not packages:
            return

        _, _, tail = os.path.splitroot(PKG._DB_REPOS_PATH)
        dest = os.path.join(rootdir, tail)

        if not os.path.exists(dest):
            self._makedirs(dest)

        ctx = TMPFSContext(self._module, dest)
        with ctx:
            if ctx.err:
                return ctx.err

            if (err := self._install(packages, rootdir)) is not None:
                return err

        return ctx.err

    def _makedirs(self, path: str) -> Error | None:
        if not self._module.check_mode:
            try:
                os.makedirs(name=path, mode=0o755, exist_ok=True)
            except OSError as err:
                return Error(str(err))

        self._changed = True
        self._diff.append(
            {
                "before_header": path,
                "before": dict(state="absent"),
                "after_header": path,
                "after": dict(state="present"),
            }
        )

    def _install(self, packages: list[str], rootdir: str) -> Error | None:
        missing_packages, err = self._filter_missing_packages(rootdir, packages)
        if err is not None:
            return err

        if not missing_packages:
            return

        return self._install_packages(missing_packages, rootdir)

    def _filter_missing_packages(
        self, rootdir: str, packages: list[str]
    ) -> tuple[list[str], Error | None]:
        if not os.path.exists(rootdir):
            return packages, None

        # Get a list of installed package names.
        cmd = [PKG._PKG_CMD, "--rootdir", rootdir, "query", "%n"]

        rc, stdout, stderr = self._module.run_command(cmd)
        if rc != os.EX_OK:
            return [], Error(stderr.strip())

        installed_packages = set(stdout.splitlines())
        if not installed_packages:
            return packages, None

        # The requested package will get through if it includes the package
        # version.
        missing_packages = list(set(packages) - installed_packages)
        missing_packages.sort()

        return missing_packages, None

    def _install_packages(self, packages: list[str], rootdir: str) -> Error | None:
        if not self._module.check_mode:
            cmd = [PKG._PKG_CMD, "--rootdir", rootdir, "install", "--yes"] + packages

            rc, _, stderr = self._module.run_command(cmd)
            if rc != os.EX_OK:
                return Error(stderr.strip())

        self._installed_packages = packages

        self._changed = True
        self._diff.extend(
            {
                "before_header": p,
                "before": dict(state="absent"),
                "after_header": p,
                "after": dict(state="present"),
            }
            for p in packages
        )


def run(module: AnsibleModule) -> tuple[dict, Error | None]:
    path = module.params["path"]
    packages = module.params["packages"]

    pkg = PKG(module)

    if (err := pkg.bootstrap(path)) is not None:
        return {}, err

    if (err := pkg.install(packages, rootdir=path)) is not None:
        return {}, err

    return (
        dict(
            path=path,
            packages=pkg.installed_packages,
            changed=pkg.changed,
            diff=pkg.diff,
        ),
        None,
    )


def main() -> None:
    module = AnsibleModule(
        argument_spec=dict(
            path=dict(type="str", required=True),
            packages=dict(type="list", required=True),
        ),
        supports_check_mode=True,
    )

    try:
        res, err = run(module)
        if err != None:
            module.fail_json(msg=str(err))
        module.exit_json(**res)
    except Exception:
        module.fail_json(msg=traceback.format_exc())


if "__main__" == __name__:
    main()
