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

Setup cluster with Ansible (`-K` option is to prompt for the `root` password on
the server, needed at the very first step to install doas(1) for further remote
management):

```
$ cd ansible/

$ ansible-playbook -K ./setup-lab.yml
```
