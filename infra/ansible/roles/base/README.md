<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Base

Base configures a FreeBSD server with following changes:

- boot: remove boot delay.
- pkg: use latest packages, add `outdated` and `vital` aliases, install NeoVim.
- syslog: migrate RC variable to /etc/rc.conf.d/syslogd and enable console
  logging.
- ntp: migrate rc-vars to /etc/rc.conf.d/ntpd and listen on local IP.
- ssh: migrate rc-vars to /etc/rc.conf.d/sshd, listen on bridge interfaces, and
  restrict allowed users to `op`.
- periodic: save daily logs.
- zfs: migrate rc-vars to /etc/rc.conf.d/zfs and enable scrubs every 2 weeks.
- pf: enable pf, pflog, and configure with bridge0

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
