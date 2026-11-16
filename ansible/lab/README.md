<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**ansible** - configure LAB cluster with Ansible

# DESCRIPTION

`ansible/` automates Lab cluster management with
[Ansible](https://docs.ansible.com).

Start with a fresh installation of FreeBSD on the server. Make sure to create
a user for Ansible to manage the host:

```
user: op
groups: op,wheel
```

Run an SSH server and install `op` user SSH keys for password-less remote
access.

Install `ansible-core` with Poetry and Ansible collections with
`ansible-galaxy` (collections go to `~/.ansible/collections`). `make -C home`
does both, or by hand from the repository root:

```
$ poetry install --no-root

$ poetry run ansible-galaxy collection install -r ansible/requirements.yml
```

Setup cluster with Ansible (`-K` option is to prompt for the `root` password on
the server, needed at the very first step to install doas(1) for further remote
management):

```
$ cd ansible/

$ poetry run ansible-playbook -K ./setup-lab.yml
```
