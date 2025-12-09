Name
====

**freebsd** - development setup under FreeBSD

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
% cat /usr/local/etc/doas.conf
permit nopass :op cmd pkg
```

Git
---

Set user name and email for commit messages:

```console
% git config --global user.name 'Samvel Khalatyan'
% git config --global user.email sn.khalatyan@gmail.com
```

Use SSH keys to sign commits. Generate a key with ssh-keygen(1). Upload the public key to GitHub. Place the public key into the environment with git(1), say `~/.ssh//id_github.pub`\):

```console
% git config --global gpg.format ssh
% git config --global user.signingkey $HOME/.ssh/id_github.pub
```

Users & Groups
--------------

Add user `op` to the `:wheel` group (don't make LDAP an authoritative source for this, don't store wheel in LDAP for security reasons. We want local system be authoritative source for root, wheel, and other widely used groups).

```console
# pw groupmod -n wheel -m op.
```
