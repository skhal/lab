<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**ansible** - configure LAB cluster with Ansible

# DESCRIPTION

Use Ansible (https://docs.ansible.com) to configure Lab cluster from
laptop over SSH.

Make sure that NUC has the following setup:

- User `op` is a member of `wheel` group.
- Running SSH server with `op` SSH key installed to let Ansible run
  commands without password prompts.

Install Ansible on the laptop:

```console
# brew install ansible
```

## Configure

Bootstrap NUC with Ansible to install python and setup user environment (`-K`
to prompt for NUC root-password):

```console
$ ansible-playbook -K ./nuc_bootstrap.yaml
```

Configure NUC without root privileges:

```console
$ ansible-playbook ./nuc.yaml
```

## Test

Test configuration with:

```console
$ ansible-playbook nuc.yml --check --diff
$ ansible-playbook nuc.yml --check --diff --tags dumpdev
```
