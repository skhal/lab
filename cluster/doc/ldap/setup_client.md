Name
====

**setup-client** - setup SSO and SSH with LDAP server

LDAP Client
===========

```console
# pkg install openldap26-client
```

Check the package message for setup instructions:

```console
# pkg info -D openldap26-client
openldap26-client-2.6.10:
On install:
The OpenLDAP client package has been successfully installed.

Edit
  /usr/local/etc/openldap/ldap.conf
to change the system-wide client defaults.

Try `man ldap.conf' and visit the OpenLDAP FAQ-O-Matic at
  http://www.OpenLDAP.org/faq/index.cgi?file=3
for more information.
```

Patch LDAP client configuration:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/ldap/client/openldap_ldap.conf.diff
# patch -lb -i /tmp/openldap_ldap.conf.diff /usr/local/etc/openldap/ldap.conf
```

Verify:

```console
# ldapsearch -x dn | grep '^dn:'
dn: dc=lab,dc=net
dn: ou=people,dc=lab,dc=net
dn: ou=groups,dc=lab,dc=net
dn: uid=op,ou=people,dc=lab,dc=net
dn: cn=op,ou=groups,dc=lab,dc=net
```

Single Sign On
==============

Name-Service Switch (NSS) dispatcher nsdispatch(3) uses NSS configuration nsswitch.conf(5) to look up various databases for information, e.g. groups, passwords, etc.

Install NSS LDAP plugin:

```console
# pkg install nss_ldap
# pkg info -D nss_ldap
nss_ldap-1.265_15:
On install:
The nss_ldap module expects to find its configuration files at the
following paths:

LDAP configuration:     /usr/local/etc/nss_ldap.conf
LDAP secret (optional): /usr/local/etc/nss_ldap.secret
```

Patch NSS LDAP configuration:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/ldap/client/nss_ldap.conf.diff
# patch -lb -i /tmp/nss_ldap.conf.diff /usr/local/etc/nss_ldap.conf
```

Change NSS configuration to pull Single Sign On (SSO) groups and passwords from files then LDAP. It is important to give preference to give local users and groups preference to act as the authoritative source and consult LDAP only if local entry is not found, i.e., the order should be `files ldap`:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/ldap/client/nsswitch.conf.diff
# patch -lb -i /tmp/nsswitch.conf.diff /etc/nsswitch.conf
```

Verify:

```console
# getent group op
op:*:1000:op
# getent passwd op
op:*:1000:1000:Operator:/home/op:/bin/tcsh
```

Secure Shell
============

Use Pluggable Authentication Modules (PAM) to configure SSH with LDAP:

```console
# pkg install pam_ldap pam_mkhomedir
# pkg info -D pam_ldap
pam_ldap-186_2:
On install:
Edit /usr/local/etc/ldap.conf in order to use this module.  Then
create a /usr/local/etc/pam.d/ldap with a line similar to the following:

login	auth	sufficient	/usr/local/lib/pam_ldap.so
```

-	`pam_ldap` provides LDAP authentication.
-	`pam_mkhomedir` creates user home folder if one does not exist.

Keep in mind that PAM LDAP uses `/usr/local/etc/ldap.conf` for configuration, compared to OpenLDAP client configuration at `/usr/local/etc/openldap/ldap.conf`, i.e., the same filename but different paths. Patch it to point PAM to LDAP server:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/ldap/client/ldap.conf.diff
# patch -lb -i /tmp/ldap.conf.diff /usr/local/etc/ldap.conf
```

PAM policies reside under `/etc/pam.d` in pam.conf(5) format (or `/usr/local/etc/pam.d` for installed packages). Patch SSH policy to use LDAP:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/ldap/client/sshd.diff
# patch -lb -i /tmp/sshd.diff /etc/pam.d/sshd
```

Verify:

```console
# ssh op@dev.lab.net whoami
(op@dev.lab.net) Password:
op
```
