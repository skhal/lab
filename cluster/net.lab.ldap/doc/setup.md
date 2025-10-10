# NAME

**setup** - setup of `ldap.nuc.lab.net`

# DESCRIPTION

> [!WARNING]
> The setup is minimal, without traffic encryption. SSL/TLS encryption is coming
> shortly.

## Client setup

OpenLDAP configuration is found at `/usr/local/etc/openldap/ldap.conf`. It
stores default values for LDAP clients. We'll set the default base to be
`dc=lab,dc=net` and use local LDAP server at IP `192.168.1.90`.

```console
# cat /usr/local/etc/openldap/ldap.conf
#
# LDAP Defaults
#

# See ldap.conf(5) for details
# This file should be world readable but not world writable.

BASE    dc=lab,dc=net
URI     ldap://192.168.1.90

#SIZELIMIT      12
#TIMELIMIT      15
#DEREF          never
```

Use password we have set in the slapd(8) configuration under `olcRootPW` for
`dc=lab,dc=net` database to create a Directory Information Tree (DIT):

```console
# ldapadd -H ldap://192.168.1.90 -xW -D 'cn=op,dc=lab,dc=net' -f ./openldap/lab.net/1.create_dit.ldif
```

> [!TIP]
> Use the following tcsh(1) alias to speed up search commands:<br/>
> `alias ldapsearch /usr/local/bin/ldapsearch -H ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi -x`

Verify the changes:

```console
# ldapsearch dn | grep '^dn:'
dn: dc=lab,dc=net
dn: ou=people,dc=lab,dc=net
dn: ou=groups,dc=lab,dc=net
```

## Manage LDAP

Any client can update LDAP server contents with instructions in LDIF format:

```console
% ldapadd -x -W -D 'cn=op,dc=lab,dc=net' -f /tmp/foo.ldif
% ldapmodify -x -W -D 'cn=op,dc=lab,dc=net' -f /tmp/foo.ldif
```

Change user password:

```console
% ldappasswd -x -W -D 'cn=op,dc=lab,dc=net' -S 'uid=foo,ou=people,dc=lab,dc=net'
```
