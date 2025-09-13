# NAME

**setup** - basic setup of `dev.nuc.lab.net`


# USERS & GROUPS

Users:

  * `skhalatyan` is a software engineer

    ```console
    # pw groupadd \
        -g 10001 \
        -n skhalatyan
    # pw useradd \
        -c 'Software Engineer' \
        -d /home/skhalatyan \
        -g skhalatyan \
        -m \
        -n skhalatyan \
        -s /bin/tcsh \
        -u 10001 \
        -w no
    # passwd -l skhalatyan
    ```

  * `op` is an operator of the host

    ```console
    # pw groupadd \
        -g 1000 \
        -n op
    # pw useradd \
        -c 'System Operator' \
        -d /home/op \
        -g op \
        -G wheel \
        -m \
        -n op \
        -s /bin/tcsh \
        -u 1000 \
        -w no
    # passwd -l op
    ```


Groups:

  * `ssh` lists users that are allowed to SSH into the host

    > [!WARNING]
    > the command adds **existing** users to the group with `-M ...` flag.

    ```console
    # pw groupadd \
        -g 1001 \
        -n ssh \
        -M op,skhalatyan
    ```

# SERVICES

## SSH

```console
# sysrc -f /etc/rc.conf.d/sshd sshd_enable=yes
# sysrc -f /etc/rc.conf.d/sshd sshd_flags="-o ListenAddress=192.168.1.110 -o AllowGroups=ssh"
# service sshd restart
# sockstat -4 | grep sshd
root     sshd       18351 7   tcp4   192.168.1.110:22      *:*
```


# Applications

Switch from Quarterly to Latest packages:

```console
# mkdir -p /usr/local/etc/pkg/repos
# cat <<eof > /usr/local/etc/pkg/repos/FreeBSD.conf 
FreeBSD: { url: "pkg+http://pkg.FreeBSD.org/${ABI}/latest" }
eof
# pkg update -f
```

Installed packages:

```console
% pkg prime-list
bazel
doas
git
go125
llvm
pkg
rsync
tmux
tree
vim
```

The following packages keep installation under own folder instead of injecting
files into appropriate locations of `/usr/local/`:

  * `go125` goes into `/usr/local/go125`
  * `llvm` goes into `/usr/local/llvm19`


# Other

## Date & time

The date & time might be using UTC:

```console
% date
Fri Sep  5 03:12:45 UTC 2025
```

Run `tzsetup(8)` to change timezone.

```console
% date
Thu Sep  4 22:19:54 CDT 2025
```
