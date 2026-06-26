<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Bootstrap

Bootstrap configures FreeBSD server for Ansible automation.

This role installs Python3 and doas(1).

## Requirements

The FreeBSD server must have:

- a running SSH server.

- an `op` user with `wheel` group membership, SSH access, and SSH keys
  installed.

## Example Playbook

```yaml
# file: bootstrap.yml
---
- name: Bootstrap

  hosts: example
  gather_facts: false

  become: true
  become_method: ansible.builtin.su

  roles:
    - role: bootstrap
```

Run the playbook with become-pass:

```console
$ ansible-playbook -K ./bootstrap.yml
```

## License

BSD-3-Clause

## Author Information

[Samvel Khalatyan](https://github.com/skhal)
