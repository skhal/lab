Name
====

**freebsd** - development setup under FreeBSD

Bootstrap
=========

```console
% pkg install FreeBSD-ssh FreeBSD-bmake doas
```

Configure ssh:

```
# sysrc -f /etc/rc.conf.d/sshd sshd_enable=yes
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/dev/data/sshd_config.diff
# patch -lb -i /tmp/sshd_config.diff /etc/ssh/sshd_config
```

Install
=======

The repository comes with FreeBSD make files to install software:

```console
% make install
```

Important [packages](https://github.com/skhal/lab/blob/63436f8239a1447b062f5d64dbafd0349642bc58/home/pkg/packages.txt):

-	`go125` for Go
-	`llvm21` for Vim LSP with `clangd`
-	`tmux` terminal multiplexer
-	`universal-ctags` for Vim tagbar with C++

It also installs statically linked Go [binaries](https://github.com/skhal/lab/blob/63436f8239a1447b062f5d64dbafd0349642bc58/home/go/bin/packages.txt) (static linking makes it work in Linux jails with shared home folder):

-	`ibazel` interactive Bazel
-	`gopls` Go LSP for Vim
-	various pre-commit checks

Applications
============

Doas
----

Give privileges to the members of `:op` group install packages:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/dev/data/doas.conf
# mv /tmp/doas.conf /usr/local/etc/
```

Fonts
-----

`webfonts` package installs fonts. Re-generate font cache files:

```console
# fc-cache -f
```

Validate - should be non-empty list of installed fonts:

```console
% fc-list | wc -l
  30
```

Git
---

Set user name and email for commit messages:

```console
% git config --global user.name 'John Doe'
% git config --global user.email john.doe@example.com
```

Use SSH keys to sign commits ([ref](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)). Generate a key with ssh-keygen(1). Upload the public key to GitHub. Place the public key into the environment with git(1), say `~/.ssh//id_github.pub`\):

```console
% git config --global user.signingkey ~/.ssh/id_github.pub
```

Packages
--------

Use latest package for development:

```console
% cat /usr/local/etc/pkg/repos/FreeBSD.conf
FreeBSD-ports: {
  url: "pkg+https://pkg.FreeBSD.org/${ABI}/latest",
}
```

Users & Groups
--------------

Add user `op` to the `:wheel` group (don't make LDAP an authoritative source for this, don't store wheel in LDAP for security reasons. We want local system be authoritative source for root, wheel, and other widely used groups).

```console
# pw groupmod -n wheel -m op
```
