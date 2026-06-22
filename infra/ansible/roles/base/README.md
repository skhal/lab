<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Base

Base configures a FreeBSD server.

## Requirements

The FreeBSD server must be bootstrapped for Ansible:

- an SSH server running.
- has an `op` user with `wheel` group membership and SSH keys installed.

## Role Variables

All variables are internal.

## Example Playbook

```
- name: Configure example server
  hosts: example
  gather_facts: false
  become: true
  become_method: community.general.doas

  roles:
    - role: base
```

## License

BSD-3-Clause

## Author Information

[Samvel Khalatyan](https://github.com/skhal)
