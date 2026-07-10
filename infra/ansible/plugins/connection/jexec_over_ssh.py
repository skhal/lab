# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import dataclasses
import json
import os
from contextlib import contextmanager
from dataclasses import dataclass
from functools import wraps
from typing import Callable, Optional, override

import yaml
from ansible.errors import AnsibleError
from ansible.module_utils.common.text.converters import to_text
from ansible.plugins.connection import ssh
from ansible.plugins.connection.ssh import Connection as SSHConnection

DOCUMENTATION = r"""
name: jexec_over_ssh
short_description: Run tasks using jexec on the remote host.
description:
    - Run commands or put/fetch files to a jail on the remote host using jexec.
options:
    jail_name:
        description:
            - Jail name to connect to.
        type: str
        required: true
        vars:
            - name: ansible_jexec_over_ssh_jail_name
author: Samvel Khalatyan
"""

EXAMPLES = r"""
# SSH to bar.example.com and jexec to the jail foo, assuming the jail is
# registered in the inventory under foo.example.com hostname.
- name: Example
  hosts: foo.example.com
  connection: jexec_over_ssh
  vars:
    # ssh to bar.example.ecom
    ansible_ssh_host: bar.example.com
    # jexec to foo jail
    ansible_jexec_over_ssh_jail_name: foo
  gather_facts: false
  # elevate privileges to run jexec.
  become: true
  become_method: community.general.doas
  tasks:
    - name: Ping the jail
      ansible.builtin.ping:

# Same but using Ansible CLI:
#   ansible foo.example.com \
#       -b \
#       --become-method=community.general.doas \
#       -c jexec_over_ssh \
#       -e ansible_jexec_over_ssh_jail_name=foo \
#       -e ansible_ssh_host=bar.example.com \
#       -m ping
"""

_DOC_OPTIONS_KEY = "options"


def _merge_doc_options(src: str, dst: str) -> str:
    """Merge documentation options from src to dst where src options have
    overwrite priority.
    """
    doc = yaml.safe_load(dst) or {}
    opts = doc.get(_DOC_OPTIONS_KEY, {})
    opts.update((yaml.safe_load(src) or {}).get(_DOC_OPTIONS_KEY, {}))
    doc[_DOC_OPTIONS_KEY] = opts
    return yaml.safe_dump(doc)


globals().update(DOCUMENTATION=_merge_doc_options(ssh.DOCUMENTATION, DOCUMENTATION))

# keep-sorted start
_CP_CMD = "/bin/cp"
_JEXEC_CMD = "/usr/sbin/jexec"
_JLS_CMD = "/usr/sbin/jls"
_MKTEMP_CMD = "/usr/bin/mktemp"
_RM_CMD = "/bin/rm"
# keep-sorted end


@dataclass
class Jail:
    name: str
    path: str


class JailNotFoundError(AnsibleError):
    """JailNotFoundError means specified jail is not found."""

    def __init__(self, name: str):
        super().__init__(f"jail {name} not found")


def as_jail_list(d: dict) -> dict | Jail:
    """Remove wrapper json header and parse jails.

    jls(8) JSON output is:
        {
          "__version": "2",
          "jail-information": {
            "jail": [
              {
                "jid": 123
                "name": "foo",
                "path": "/jail/foo"
              }
            ]
          }
        }

    This function parses every jail as Jail structure and removes the wrapper
    layer at "jail-information". The returned data are:

        dict(jail=[
            Jail(name="foo", path="/jail/foo"),
        ])
    """
    jail_information_key = "jail-information"
    if jail_information_key in d:
        return d[jail_information_key]
    if "jid" not in d:
        return d
    # jls output parameters must match Jail field names
    keys = {f.name for f in dataclasses.fields(Jail)}
    return Jail(**{k: d[k] for k in d.keys() & keys})


class Connection(SSHConnection):
    """Remote BSD jail connection."""

    _OPT_JAIL_NAME = "jail_name"

    transport = "jexec_over_ssh"  # plugin identifier
    has_pipelining = True  # use persistent connection

    def __init__(self, *args, **kwargs) -> None:
        super().__init__(*args, **kwargs)
        self._jail = None
        self._remote_host = None

    @property
    def jail(self) -> Jail:
        if not self._jail:
            self._jail = self._init_jail()
        return self._jail

    @property
    def remote_host(self) -> str:
        if not self._remote_host:
            self._remote_host = (
                self.get_option("host") or self._play_context.remote_addr
            )
        return self._remote_host

    def _init_jail(self) -> Jail:
        """Pull jail properties from the remote host using jls(8)."""
        jname = self.get_option(Connection._OPT_JAIL_NAME)
        cmd = [_JLS_CMD, "--libxo=json", "-j", jname, "jid", "name", "path"]
        rc, out, _ = super().exec_command(" ".join(cmd))
        if rc != os.EX_OK:
            raise JailNotFoundError(jname)
        jail = self._parse_jls_output(to_text(out.strip()), jname)
        self._display.vvv(f"{self.transport}: validated jail {jail}")
        return jail

    @staticmethod
    def _parse_jls_output(out: str, jail_name: str) -> Jail:
        jail_list = json.loads(out, object_hook=as_jail_list)
        jails = {j.name: j for j in jail_list.get("jail", [])}
        if jail_name not in jails:
            raise JailNotFoundError(jail_name)
        return jails[jail_name]

    @override
    def exec_command(
        self, cmd: str, in_data=None, sudoable=True
    ) -> tuple[int, bytes, bytes]:
        """Execute a remote command cmd with optional stdin in_data and
        escalated privileges when sudoable is set to True.
        """
        self._display.vvv(
            f"{self.transport}: original command: {cmd}",
            host=f"{self.remote_host} {self.jail.name}",
        )
        jexec_cmd = [
            _JEXEC_CMD,
            "-l",  # run in clean environment
            self.jail.name,
            cmd,
        ]
        return self._exec_command(jexec_cmd, in_data, sudoable)

    def _exec_command(
        self, command: list[str], in_data=None, sudoable=True
    ) -> tuple[int, bytes, bytes]:
        cmd = " ".join(command)
        if self.become:
            cmd = self.become.build_become_command(cmd, self._shell)
        self._display.vvv(
            f"{self.transport}: {cmd}",
            host=f"{self.remote_host} {self.jail.name}",
        )
        return super().exec_command(cmd, in_data=in_data, sudoable=sudoable)

    @override
    def put_file(self, in_path: str, out_path: str) -> tuple[int, bytes, bytes]:
        """Transfer a local in_path file to the remote out_path file."""
        with self._remote_mktemp() as tmp_file:
            rc, out, err = super().put_file(in_path, tmp_file)
            if rc != os.EX_OK:
                return rc, out, err
            cp_cmd = [
                _CP_CMD,
                "-f",
                "-p",  # preserve attributes
                tmp_file,
                os.path.join(self.jail.path, out_path),
            ]
            return self._exec_command(cp_cmd)

    @override
    def fetch_file(self, in_path: str, out_path: str) -> tuple[int, bytes, bytes]:
        """Fetch a remote in_path file to the local out_path file."""
        with self._remote_mktemp() as tmp_file:
            cp_cmd = [
                _CP_CMD,
                "-f",
                "-p",  # preserve attributes
                os.path.join(self.jail.path, in_path),
                tmp_file,
            ]
            rc, out, err = self._exec_command(cp_cmd)
            if rc != os.EX_OK:
                return rc, out, err
            return super().fetch_file(tmp_file, out_path)

    @contextmanager
    def _remote_mktemp(self):
        """Creates a temporary file on the remote host. It removes the file upon
        exit from the context."""
        cmd = [_MKTEMP_CMD, "--quiet"]
        rc, out, _ = super().exec_command(" ".join(cmd))
        if rc != os.EX_OK:
            raise AnsibleError(f"failed to create temporary file on {self.remote_host}")
        tmp_file = to_text(out.strip())
        try:
            yield tmp_file
        finally:
            cmd = [_RM_CMD, "-f", tmp_file]
            super().exec_command(" ".join(cmd))
