# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import override
import yaml

from ansible.plugins.connection import ssh
from ansible.plugins.connection.ssh import Connection as SSHConnection
from ansible.utils.display import Display

DOCUMENTATION = r"""
name: jexec_over_ssh
short_description: Run tasks using jexec on the remote host.
description:
    - Run commands or put/fetch files to a jail on the remote host using jexec.
author: Samvel Khalatyan
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

display = Display()


class Connection(SSHConnection):
    """Remote BSD jail connection."""

    transport = "jexec_over_ssh"  # plugin identifier
    has_pipelining = True  # use persistent connection

    def __init__(self, *args, **kwargs) -> None:
        super(Connection, self).__init__(*args, **kwargs)

    @override
    def _connect(self) -> Connection:
        """Establish a connection to the remote host."""
        super(Connection, self)._connect()
        display.vvv("JEXEC-OVER-SSH: connected", host=self._play_context.remote_addr)
        return self

    @override
    def exec_command(
        self, cmd, in_data=None, sudoable=True
    ) -> tuple[int, bytes, bytes]:
        """Execute a remote command cmd with optional stdin in_data and
        escalated privileges when sudoable is set to True.
        """
        return super(Connection, self).exec_command(
            cmd, in_data=in_data, sudoable=sudoable
        )

    @override
    def put_file(self, in_path: str, out_path: str) -> tuple[int, bytes, bytes]:
        """Transfer a local in_path file to the remote out_path file."""
        return super(Connection, self).put_file(in_path, out_path)

    @override
    def fetch_file(self, in_path: str, out_path: str) -> tuple[int, bytes, bytes]:
        """Fetch a remote in_path file to the local out_path file."""
        return super(Connection, self).fetch_file(in_path, out_path)

    @override
    def close(self) -> None:
        """Close the remote connection."""
        super(Connection, self).close()
