NAME
====

**setup** - setup Ubuntu 22.04 Jammy jail

DESCRIPTION
===========

Clang
-----

Setup LLVM apt(1) sources:

```console
% doas jexec jammy chroot /compat/jammy su -l
root@jammy:~# wget -qO /etc/apt/trusted.gpg.d/apt.llvm.org.asc https://apt.llvm.org/llvm-snapshot.gpg.key
root@jammy:~# cat /etc/apt/sources.list.d/llvm.list
# Reference: https://apt.llvm.org
deb http://apt.llvm.org/jammy/ llvm-toolchain-jammy-19 main
deb-src http://apt.llvm.org/jammy/ llvm-toolchain-jammy-19 main
```

Install clang(1) from LLVM 19:

```console
% doas jexec jammy chroot /compat/jammy apt install clang-19
```
