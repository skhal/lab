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
% doas jexec jammy chroot /compat/jammy cat /etc/apt/sources.list
deb http://archive.ubuntu.com/ubuntu jammy main universe restricted multiverse
deb http://security.ubuntu.com/ubuntu/ jammy-security universe multiverse restricted main
deb http://archive.ubuntu.com/ubuntu jammy-backports universe multiverse restricted main
deb http://archive.ubuntu.com/ubuntu jammy-updates universe multiverse restricted main
```

Upgrade the user land:

```console
% doas jexec jammy chroot /compat/jammy apt update
% doas jexec jammy chroot /compat/jammy apt upgrade
```

## Tools

Add tcsh(1) for LDAP users. Simplify configuration edits with vim(1):

```console
% doas jexec jammy chroot /compat/jammy apt install vim tcsh
```

## Locale

```console
% doas jexec jammy chroot /compat/jammy locale-gen C.UTF-8
Generating locales (this might take a while)...
  C.UTF-8... done
Generation complete.
% doas jexec jammy chroot /compat/jammy \
    dpkg-reconfigure --frontend noninteractive locales
Generating locales (this might take a while)...
  C.UTF-8... done
Generation complete.
```

## Timezone

```console
% doas jexec jammy chroot /compat/jammy \
    ln -fs /usr/share/zoneinfo/America/Chicago /etc/localtime
% doas jexec jammy chroot /compat/jammy \
    dpkg-reconfigure --frontend noninteractive tzdata

Current default time zone: 'America/Chicago'
Local time is now:      Sat Oct  4 09:10:13 CDT 2025.
Universal Time is now:  Sat Oct  4 14:10:13 UTC 2025.
```

Verify:

```console
% doas jexec jammy chroot /compat/jammy date
Sat Oct  4 09:10:38 CDT 2025
```

## Name Service Switch

Confirm LDAP server is reachable. Install ldapsearch(1):

```console
% doas jexec jammy chroot /compat/jammy apt install ldap-utils
```

Verify:

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

Configure with LDAP server:
  * IP: 10.0.1.90
  * Base: `dc=lab,dc=net`
  * Services: `passwd` and `group`

Disable TLS for now:

```console
% doas jexec jammy chroot /compat/jammy \
    sed -i 's/^tls_cacertfile/#\0/' /etc/nslcd.conf
```

Start nslcd(8) service in jail's `/etc/rc.local` using Linux service(8):

```console
% doas jexec jammy cat /etc/rc.local
chroot /compat/jammy /usr/sbin/service nslcd start
% doas jexec jammy cat /etc/rc.shutdown.local
chroot /compat/jammy /usr/sbin/service nslcd stop
```

> [!IMPORTANT]
> Make sure there is a single line to start each service in `/etc/rc*local`.
> Also, start/stop nslcd(8) before/after other services, especially ssh(8).

Restart the jail to start nslcd(8) service.

```console
% doas jexec jammy chroot /compat/jammy service --status-all |& grep nslcd
 [ + ]  nslcd
% doas jexec jammy ps -Af |& grep nslcd
68288  -  IJ   0:00.00 /usr/sbin/nslcd
% doas jexec jammy sockstat -lu | grep nslcd
106      nslcd      68288 5   stream /var/run/nslcd/socket
```

Verify work:

```
% doas jexec jammy chroot /compat/jammy getent passwd op
op:*:1000:1000:Operator:/home/op:/bin/tcsh
```

## PAM

```console
% doas jexec jammy chroot /compat/jammy apt install libpam-ldapd
```

> [!NOTE]
> The package was probably installed along `libnss-ldap`.

## SSH

```console
% doas jexec jammy chroot /compat/jammy apt install openssh-server
```

Configure:

```console
% doas jexec jammy chroot /compat/jammy su -l
root@jammy:~# cat <<eof >/etc/ssh/sshd_config.d/lab.net.conf
> ListenAddress=192.168.1.111
> AllowGroups ssh
> eof
root@jammy:~# exit
```

Start the service:

```console
% doas jexec jammy cat /etc/rc.local
chroot /compat/jammy /usr/sbin/service nslcd start
chroot /compat/jammy /usr/sbin/service ssh start
% doas jexec jammy cat /etc/rc.shutdown.local
chroot /compat/jammy /usr/sbin/service ssh start
chroot /compat/jammy /usr/sbin/service nslcd stop
```

> [!IMPORTANT]
> Make sure there is a single line to start each service in `/etc/rc*local`.

Restart the jail to start sshd(8) service.

```console
% doas jexec jammy chroot /compat/jammy service --status-all |& grep ssh
 [ + ]  ssh
% doas jexec jammy ps -Af | grep ssh
97160  -  IsJ  0:00.00 /usr/sbin/sshd
% doas jexec jammy sockstat -l4 | grep ssh
root     sshd       97160 3   tcp4   192.168.1.111:22      *:*
```

Verify work:

> [!WARNING]
> SSH may refuse to accept credentials if user's shell is not installed, e.g.
> confirm that tcsh(1) is available.

```console
% ssh op@jammy.lab.net hostname
op@jammy.lab.net's password:
jammy.lab.net
```
