# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""A Poetry application plugin for `poetry patch` command.

Poetry loads plugins at startup and adds plugin commands to the CLI.
"""

from poetry.console.commands.command import Command
from poetry.plugins import ApplicationPlugin

from poetry_plugin_patch.command import PatchCommand


class PatchApplicationPlugin(ApplicationPlugin):
    """The plugin adds `poetry patch` command."""

    @property
    def commands(self) -> list[type[Command]]:
        """A list of new commands: PatchCommand."""

        return [PatchCommand]
