<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**ubuntu** -- Ubuntu development environment

# DESCRIPTION

## Networking

```console
$ sudo apt install sockstat
```

## Secure Shell

```console
$ sudo apt install openssh-server
$ sudo systemctl enable sshd
```

## Tools

```console
$ sudo apt install neovim tmux
```

## User

Ubuntu 26.04 installation requires one to create a user with sudo(1) privieleges
to operate the system. There is a high chance the user will have 1000 for user
and group id.

Create another user with a different uid/gid as necessary, e.g. `admin`:

```console
$ sudo addgroup --gid 10000 admin
$ sudo adduser --gid 10000 --uid 10000 admin
$ sudo usermod -aG adm,sudo admin
```

Ubuntu uses `adm` group to let system administrator monitor the system. Many
logs, configurations belong to this group. It is a good idea to add the new
user to `adm`.

Finally, remove the old user, e.g. `op`:

```console
$ sudo deluser --remove-home op
```

Set tcsh(1) for default shell:

```console
$ sudo apt install tcsh
$ chsh -s /usr/bin/tcsh
```
