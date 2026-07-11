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
    jail_user:
        description:
            - Jail user to run commands as using jexec(8) -U flag.
        type: str
        vars:
            - name: ansible_user
            - name: ansible_jexec_over_ssh_jail_user
    escalator:
        description:
            - Method used to escalate jexec(8) privileges.
        type: str
        choices: [doas, sudo]
        vars:
            - name: ansible_jexec_over_ssh_escalator
author: Samvel Khalatyan
"""

EXAMPLES = r"""
# Setup: a hostname foo.example.com runs a jail bar with hostname
# bar.example.com.

# Run ping module inside the jail using `doas jexec -U {{ ansible_user }}`
- name: Example
  hosts: bar.example.com
  connection: jexec_over_ssh
  vars:
    ansible_ssh_host: foo.example.com
    ansible_jexec_over_ssh_jail_name: bar
    ansible_jexec_over_ssh_escalator: doas
  tasks:
    - name: Ping the jail
      ansible.builtin.ping:

# Run ping module inside the jail with operator-user using
# `doas jexec -U operator`
- name: Example
  hosts: bar.example.com
  connection: jexec_over_ssh
  vars:
    ansible_ssh_host: foo.example.com
    ansible_jexec_over_ssh_jail_name: bar
    ansible_jexec_over_ssh_jail_user: operator
    ansible_jexec_over_ssh_escalator: doas
  tasks:
    - name: Ping the jail
      ansible.builtin.ping:

# Same but using Ansible CLI:
#   ansible foo.example.com \
#       -c jexec_over_ssh \
#       -e ansible_ssh_host=foo.example.com \
#       -e ansible_jexec_over_ssh_jail_name=bar \
#       -e ansible_jexec_over_ssh_escalator=doas \
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
_CP_CMD = "cp"
_DOAS_CMD = "doas"
_JEXEC_CMD = "jexec"
_JLS_CMD = "jls"
_MKTEMP_CMD = "mktemp"
_RM_CMD = "rm"
_SUDO_CMD = "sudo"
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


class Escalator:
    def escalate(self, cmd: list[str]) -> list[str]:
        return cmd


class DoasEscalator(Escalator):
    @override
    def escalate(self, cmd: list[str]) -> list[str]:
        return [
            _DOAS_CMD,
            "-n",  # non-interactive
            *cmd,
        ]


class SudoEscalator(Escalator):
    @override
    def escalate(self, cmd: list[str]) -> list[str]:
        return [
            _SUDO_CMD,
            "-n",  # non-interactive
            *cmd,
        ]


class Connection(SSHConnection):
    """Remote BSD jail connection."""

    # keep-sorted start
    _OPT_ESCALATOR = "escalator"
    _OPT_JAIL_NAME = "jail_name"
    _OPT_JAIL_USER = "jail_user"
    # keep-sorted end

    transport = "jexec_over_ssh"  # plugin identifier
    has_pipelining = True  # use persistent connection

    def __init__(self, *args, **kwargs) -> None:
        super().__init__(*args, **kwargs)
        self._jail = None
        self._remote_host = None
        self._escalator = None
        self._jexec_args = None

    @property
    def escalator(self) -> Escalator:
        if not self._escalator:
            self._escalator = self._init_escalator()
        return self._escalator

    def _init_escalator(self) -> Escalator:
        escalator = self.get_option(Connection._OPT_ESCALATOR)
        match escalator:
            case "doas":
                return DoasEscalator()
            case "sudo":
                return SudoEscalator()
            case None:
                return Escalator()
            case _:
                raise AnsibleError(f"unsupported escalator {escalator}")

    @property
    def jexec_args(self) -> list[str]:
        if not self._jexec_args:
            self._jexec_args = self._init_jexec_args()
        return self._jexec_args

    def _init_jexec_args(self) -> list[str]:
        args = [
            "-l",  # run in clean environment
        ]
        user = self.get_option(Connection._OPT_JAIL_USER)
        if user and user != "root":
            args.extend(["-U", user])
        return args

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
            *self.jexec_args,
            self.jail.name,
            cmd,
        ]
        jexec_cmd = self.escalator.escalate(jexec_cmd)
        return self._exec_command(jexec_cmd, in_data, sudoable)

    def _exec_command(
        self, command: list[str], in_data=None, sudoable=True
    ) -> tuple[int, bytes, bytes]:
        cmd = " ".join(command)
        self._display.vvv(
            f"{self.transport}: command: {cmd}",
            host=f"{self.remote_host} {self.jail.name}",
        )
        return super().exec_command(cmd, in_data=in_data, sudoable=sudoable)

    @override
    def put_file(self, in_path: str, out_path: str) -> tuple[int, bytes, bytes]:
        """Transfer a local in_path file to the remote out_path file."""
        self._display.vvv(
            f"{self.transport}: put {in_path} to {out_path}",
            host=f"{self.remote_host} {self.jail.name}",
        )
        with self._remote_mktemp() as tmp_file:
            rc, out, err = super().put_file(in_path, tmp_file)
            if rc != os.EX_OK:
                return rc, out, err
            dst = self._to_jail_path(out_path)
            return self._remote_copy(tmp_file, dst)

    def _to_jail_path(self, path: str) -> str:
        # os.path.join() drops arguments before out_path if the latter one
        # is an absolute path. Strip of leading slash from out_path to join
        # it to the jail path.
        if os.path.isabs(path):
            path = os.path.splitroot(path)[2]
        return os.path.join(self.jail.path, path)

    @override
    def fetch_file(self, in_path: str, out_path: str) -> tuple[int, bytes, bytes]:
        """Fetch a remote in_path file to the local out_path file."""
        self._display.vvv(
            f"{self.transport}: fetch {in_path} to {out_path}",
            host=f"{self.remote_host} {self.jail.name}",
        )
        with self._remote_mktemp() as tmp_file:
            src = self._to_jail_path(in_path)
            rc, out, err = self._remote_copy(src, tmp_file)
            if rc != os.EX_OK:
                return rc, out, err
            return super().fetch_file(tmp_file, out_path)

    def _remote_copy(self, src: str, dst: str) -> tuple[int, bytes, bytes]:
        """Copy src file to dst file on the remote host."""
        cp_cmd = [
            _CP_CMD,
            "-f",
            src,
            dst,
        ]
        return self._exec_command(cp_cmd)

    @contextmanager
    def _remote_mktemp(self):
        """Creates a temporary file on the remote host. It removes the file upon
        exit from the context."""
        cmd = [_MKTEMP_CMD, "--quiet", "-t", f"ansible.{self.transport}"]
        rc, out, _ = super().exec_command(" ".join(cmd))
        if rc != os.EX_OK:
            raise AnsibleError(f"failed to create temporary file on {self.remote_host}")
        tmp_file = to_text(out.strip())
        try:
            yield tmp_file
        finally:
            cmd = [_RM_CMD, "-f", tmp_file]
            super().exec_command(" ".join(cmd))
