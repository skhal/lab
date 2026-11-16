# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from __future__ import annotations

DOCUMENTATION = r"""
module: pkg_bootstrap
short_description: Install packages into empty, non-running jail
description:
    - Copy canonical pkg-repo(8) keys from the server into the jail.
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
  pkg_bootstrap:
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

import filecmp
import os
import shutil
import traceback

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.lab.errors import Error


class PKG:
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

        # it is expected that the source has keys path. pkg(8) can't bootstrap
        # without keys. Fail loud if that's the case
        if not os.path.isdir(PKG._KEYS_PATH):
            return Error(f"pkg keys directory not found: {PKG._KEYS_PATH}")

        try:
            files = self._stale_files(PKG._KEYS_PATH, dest)
        except OSError as err:
            return Error(str(err))
        if not files:
            return

        if not self._module.check_mode:
            try:
                for src, dst, _ in files:
                    os.makedirs(os.path.dirname(dst), mode=0o755, exist_ok=True)
                    shutil.copy2(src, dst)
            except OSError as err:
                return Error(str(err))

        self._changed = True
        self._diff.extend(
            {
                "before_header": dst,
                "before": dict(state=state),
                "after_header": dst,
                "after": dict(state="present"),
            }
            for _, dst, state in files
        )

    @staticmethod
    def _stale_files(
        src_root: str, dest_root: str
    ) -> list[tuple[str, str, str]]:
        """Lists files that are missing or differ in dest_root.

        Returns tuples of (source, destination, state), where state is the
        current state of the destination: absent or outdated.
        """

        files = []

        def onerror(err: OSError):
            raise err

        for dirpath, _, filenames in os.walk(src_root, onerror=onerror):
            relpath = os.path.relpath(dirpath, src_root)
            for name in sorted(filenames):
                src = os.path.join(dirpath, name)
                dst = os.path.normpath(os.path.join(dest_root, relpath, name))
                if not os.path.exists(dst):
                    files.append((src, dst, "absent"))
                elif not filecmp.cmp(src, dst, shallow=False):
                    files.append((src, dst, "outdated"))

        return files

    def install(self, packages: list[str], *, rootdir: str) -> Error | None:
        """install listed packaged"""

        if not packages:
            return

        return self._install(packages, rootdir)

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

    def _install_packages(
        self, packages: list[str], rootdir: str
    ) -> Error | None:
        if not self._module.check_mode:
            cmd = [
                PKG._PKG_CMD,
                "--rootdir",
                rootdir,
                "install",
                "--yes",
            ] + packages

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
