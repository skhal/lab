Name
====

**linux** - development setup under Linux

Description
===========

The instruction assume to be run in a Linux jail under user with sudo(1) privileges.

Tools
-----

```console
% sudo apt install libdigest-sha-perl tmux
```

-	`libdigest-sha-perl` for `shasum`

Bazel
-----

Use [Bazel](https://bazel.build) to build C++ code.

**Note**: Bazel recommends Bazelisk under Ubuntu, MacOS, and Windows. Unfortunately, the tool reports FreeBSD as unsupported OS even when run in a Linux jail.

Use Bazel's apt repository ([ref](https://bazel.build/install/ubuntu#install-on-ubuntu)\):

```console
% sudo apt install apt-transport-https curl gnupg -y
% curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor >bazel-archive-keyring.gpg
% sudo mv bazel-archive-keyring.gpg /usr/share/keyrings
% echo "deb [arch=amd64 signed-by=/usr/share/keyrings/bazel-archive-keyring.gpg] https://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
```

Install:

```console
% sudo apt update
% sudo apt install bazel
```

Install git for `hedron_compile_commands` module:

```console
% sudo apt install git
```

LLVM
----

LLVM Ubuntu [instructions](https://apt.llvm.org) cover two installation methods:

-	Automatic installation script `llvm.sh`
-	Ubuntu repo

Use Ubuntu repo for better control of the system:

```console
% cat <<EOF | sudo tee /etc/apt/sources.list.d//llvm.list
# Ref: https://apt.llvm.org
deb http://apt.llvm.org/jammy/ llvm-toolchain-jammy-21 main
deb-src http://apt.llvm.org/jammy/ llvm-toolchain-jammy-21 main
EOF
```

Update the keys:

```console
% wget -qO- https://apt.llvm.org/llvm-snapshot.gpg.key | sudo tee /etc/apt/trusted.gpg.d/apt.llvm.org.asc
```

Install Clang and co:

```console
% sudo apt update
% sudo apt install clang-21 clang-tools-21 clang-21-doc libclang-common-21-dev libclang-21-dev libclang1-21 clang-format-21 python3-clang-21 clangd-21 clang-tidy-21
```

Make sure to include the following tools:

-	clang
-	clang-format
-	clangd

Go
--

Install Go for Gazelle under Bazel to generate `BUILD` files with:

```console
% bazel build :gazelle
```

Even though Gazelle can install Go, we use a local Go installation on the host due to unconventional setup, e.g. Linux in FreeBSD jailed environment. In short, Gazelle downloads Go for Linux but Go linker fails ([ref](https://github.com/skhal/lab/blob/63436f8239a1447b062f5d64dbafd0349642bc58/MODULE.bazel#L21-L37)).

Use supplied script to install Go for Linux and mix in Go linker from FreeBSD distribution:

```console
% ./go_install
```
