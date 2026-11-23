# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from poetry.console.commands.command import Command
from poetry.plugins import ApplicationPlugin
from poetry_plugin_patch.command import PatchCommand


class PatchApplicationPlugin(ApplicationPlugin):
    @property
    def commands(self) -> list[type[Command]]:
        return [PatchCommand]
