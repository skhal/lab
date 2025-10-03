# NAME

**setup** - setup Ubuntu 22.04 Jammy jail


# DESCRIPTION

## RC

Enable debugging:

```console
% doas jexec jammy sysrc rc_debug=yes
rc_debug:  -> yes
% doas jexec jammy sysrc rc_info=yes
rc_info: NO -> yes
```

## Apt

Add universe and multiverse sources to apt(8):

```console
% doas jexec jammy cat /compat/jammy/etc/apt/sources.list
deb http://archive.ubuntu.com/ubuntu jammy main universe restricted multiverse
deb http://security.ubuntu.com/ubuntu/ jammy-security universe multiverse restricted main
deb http://archive.ubuntu.com/ubuntu jammy-backports universe multiverse restricted main
deb http://archive.ubuntu.com/ubuntu jammy-updates universe multiverse restricted main
```

Upgrade the user land:

```console
% doas jexec jammy chroot /compat/jammy apt update
% doas jexec jammy chroot /compat/jammy apt upgrade --yes
```

## Shell

```console
% doas jexec jammy chroot /compat/jammy apt install tcsh
```

## Locale

```console
% doas jexec jammy chroot /compat/jammy locale-gen en_US.UTF-8
Generating locales (this might take a while)...
  en_US.UTF-8... done
Generation complete.
% doas jexec jammy chroot /compat/jammy dpkg-reconfigure --frontend noninteractive locales
Generating locales (this might take a while)...
  en_US.UTF-8... done
Generation complete.
```

## Timezone

```console
% doas jexec jammy chroot /compat/jammy dpkg-reconfigure tzdata
```

## Name Service Switch

Confirm LDAP server is reachable:

```console
% doas jexec jammy chroot /compat/jammy nc -v 10.0.1.90 389
Connection to 10.0.1.90 389 port [tcp/ldap] succeeded!
^C
% doas jexec jammy chroot /compat/jammy ldapsearch -x -H ldap://10.0.1.90 -b 'dc=lab,dc=net' -s base
# extended LDIF
#
# LDAPv3
# base <dc=lab,dc=net> with scope baseObject
# filter: (objectclass=*)
# requesting: ALL
#

# lab.net
dn: dc=lab,dc=net
dc: lab
objectClass: dcObject
objectClass: organization
description: Research & Development lab
o: Research Lab

# search result
search: 2
result: 0 Success

# numResponses: 2
# numEntries: 1
```

Install `libnss-ldapd` ([doc](https://wiki.debian.org/LDAP/NSS)).

```console
% doas jexec jammy chroot /compat/jammy apt install libnss-ldapd
```

Configure with LDAP server.

Disable TLS for now:

```console
% doas jexec jammy chroot /compat/jammy diff /etc/nslcd.conf{.orig,}
28c28
< tls_cacertfile /etc/ssl/certs/ca-certificates.crt
---
> #tls_cacertfile /etc/ssl/certs/ca-certificates.crt
```

Manage nslcd(8) service under RC using Linux service(8):

```console
% doas jexec jammy su -l
root@jammy:~ # cat <<eof >>/etc/rc.local
? chroot /compat/ubuntu /usr/sbin/service nslcd start
? eof
root@jammy:~ # cat << eof >>/etc/rc.shutdown.local
? chroot /compat/ubuntu /usr/sbin/service nslcd stop
? eof
root@jammy:~ # exit
```

> [!IMPORTANT]
> Make sure there is a single line to start each service in `/etc/rc*local`.
> Also, start/stop nslcd(8) before/after other services, especially ssh(8).

Verify it works:

```console
% doas jexec jammy chroot /compat/jammy getent passwd op
op:*:1000:1000:Operator:/home/op:/bin/tcsh
```

## PAM

```console
% doas jexec jammy chroot /compat/jammy apt install libpam-ldapd
```

## SSH

```console
% doas jexec jammy chroot /compat/jammy apt install openssh-server
```

Configure:

```console
% doas jexec jammy chroot /compat/jammy su -l
root@jammy:~# cat <<eof >/etc/ssh/sshd_config.d/lab.net.conf
ListenAddress=192.168.1.111
AllowGroups ssh
eof
root@jammy:~# exit
```

Start the service:

```console
% doas jexec jammy su -l
root@jammy:~ # cat <<eof >>/etc/rc.local
? chroot /compat/ubuntu /usr/sbin/service ssh start
? eof
root@jammy:~ # cat << eof >>/etc/rc.shutdown.local
? chroot /compat/ubuntu /usr/sbin/service ssh stop
? eof
root@jammy:~ # exit
```

> [!IMPORTANT]
> Make sure there is a single line to start each service in `/etc/rc*local`.

Verity work:

> [!WARNING]
> SSH may refuse to accept credentials if shell is not installed.

```console
% ssh op@jammy.lab.net hostname
op@jammy.lab.net's password:
jammy.lab.net
```

## Dev tools

```console
% doas jexec jammy chroot /compat/jammy apt install vim tmux git
```
