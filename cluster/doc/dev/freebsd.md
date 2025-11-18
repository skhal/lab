Name
====

**freebsd** - development setup under FreeBSD

Description
===========

The repository comes with FreeBSD make files to install software:

```console
% make install
```

Important [packages](https://github.com/skhal/lab/blob/63436f8239a1447b062f5d64dbafd0349642bc58/home/pkg/packages.txt):

-	`go125` for Go
-	`llvm21` for Vim LSP with `clangd`
-	`tmux` terminal multiplexer
-	`unversal-ctags` for Vim tagbar with C++

It also installs statically linked Go [binaries](https://github.com/skhal/lab/blob/63436f8239a1447b062f5d64dbafd0349642bc58/home/go/bin/packages.txt) (static linking makes it work in Linux jails with shared home folder):

-	`ibazel` interactive Bazel
-	`gopls` Go LSP for Vim
-	various pre-commit checks
