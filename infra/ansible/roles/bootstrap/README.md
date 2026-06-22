<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Bootstrap

Bootstrap prepares a FreeBSD server for Ansible automation. It makes sure that
the server has:

- Python3 and doas(1) installed.
- Configures `root` and `op` users with tcsh(1) shell.

## Requirements

A FreeBSD server must have:

- an SSH server running.

- an `op` user with `wheel` group membership and SSH access with SSH keys
  installed.

## Role Variables

All variables are internal.

## Example Playbook

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

```
- hosts: servers
  roles:
     - { role: username.rolename, x: 42 }
```

## License

BSD-3-Clause

## Author Information

[Samvel Khalatyan](https://github.com/skhal)
