<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Sysrc

Sysrc migrates rc-vars from /etc/rc.conf to a target file using sysrc(8)
utility.

It may send an optional notification when set.

# Example

Remove `sshd_enable` from /etc/rc.conf and set it to `true` in
/etc/rc.conf.d/sshd.

```yaml
- role: sysrc
  vars:
    sysrc_migrate:
      name: sshd_enable
      value: true
      path: /etc/rc.conf.d/sshd
```

Set `sshd_flags` and send a notification to restart SSH server:

```yaml
- role: sysrc
  vars:
    sysrc_migrate:
      name: sshd_flags
      value: "-o ListenAddress=127.0.0.1"
      path: /etc/rc.conf.d/sshd
      notify: "restart_sshd"
```
