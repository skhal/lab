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
- Python 3.11 installed.

```console
# pkg install -y python
```

Install Ansible on the laptop:

```console
# brew install ansible
```

## Configure

```console
$ ansible-playbook -K -i ./inventory ./playbook/nuc_bootstrap.yaml
$ ansible-playbook -i ./inventory ./playbook/nuc.yaml
```
