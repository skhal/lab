# NAME

**setup** - setup of `ubuntu.nuc.lab.net`

# BASICS

## Security upgrade

Add universe, restricted, and multiverse sources for `apt`:

```console
root@nuc # jexec ubuntu chroot /compat/ubuntu /bin/bash
# cat /etc/apt/sources.list
deb http://archive.ubuntu.com/ubuntu jammy main universe restricted multiverse
deb http://security.ubuntu.com/ubuntu/ jammy-security universe multiverse restricted main
deb http://archive.ubuntu.com/ubuntu jammy-backports universe multiverse restricted main
deb http://archive.ubuntu.com/ubuntu jammy-updates universe multiverse restricted main
```

Pull updates and upgrade installed applications:

```console
# apt update
# apt upgrade
```

## SSH

We'd like to run SSH server from Ubuntu universe to login into chrooted
environment. Setup ssh server in Ubuntu:

```console
root@nuc # jexec ubuntu chroot /compat/ubuntu /bin/bash
# apt install openssh-server
# echo "ListenAddress 192.168.1.111" >> /etc/ssh/sshd_config
# exit
```

Start Ubuntu ssh server via jail rc(8):

```console
root@nuc # jexec ubuntu
# tail /etc/rc*.local
==> /etc/rc.local <==
chroot /compat/ubuntu /usr/sbin/service ssh start

==> /etc/rc.shutdown.local <==
chroot /compat/ubuntu /usr/sbin/service ssh stop
# exit
```

Restart the jail and conform the ssh server is up and running:

```console
root@nuc # service jail restart ubuntu
root@nuc # jexec ubuntu
# sockstat -4
USER     COMMAND    PID   FD  PROTO  LOCAL ADDRESS         FOREIGN ADDRESS      
root     sshd       84290 3   tcp4   192.168.1.111:22      *:*
```
