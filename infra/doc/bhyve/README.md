<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

*bhyve* -- BSD hypervisor

# DESCRIPTION

FreeBSD package vm-bhyve comes with CLI vm(1) to manage virtual machines (VMs)
and a vm-service to automate the VMs.

It is important to note that the vm-service has preference in the rc-order to
the jail-service:

```console
# rcorder /etc/rc.d/* /usr/local/etc/rc.d/* 2>/dev/null | cat -n | grep '/\(vm\|jail\)$'
   154	/usr/local/etc/rc.d/vm
   172	/etc/rc.d/jail
```

A VM that depends on a jailed service, e.g. DNS or LDAP, may not work as
expected.

One of the solutions is to run bhyve in a jail that depends on required jailed
services (see `deped` param in jail(8)).

# SEE ALSO

- [Run Bhyve on the host](./hosted.md)
- [Run Bhyve in a jail](./jailed.md)
