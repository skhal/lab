# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


from poetry.console.commands.command import Command


class PatchCommand(Command):
    name = "patch"
    description = "Patch installed packages."
    help = """\
The <info>patch</info> command patches installed packages in the project
environment according to the patch declarations in the pyproject.toml.
"""

    def handle(self) -> int:
        # TODO(https://github.com/skhal/lab/issues/693): implement the patch
        return 0
