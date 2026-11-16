# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""Pytest configuration for the lab Ansible tree."""

import pathlib

import ansible.module_utils

# Ansible resolves ansible.module_utils.lab.* from ./module_utils at runtime
# (see module_utils in ansible.cfg). Mirror that for pytest.
ansible.module_utils.__path__.append(
    str(pathlib.Path(__file__).parent / "module_utils")
)
