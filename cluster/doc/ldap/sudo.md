Name
====

**sudo** - add sudo(8) schema to LDAP server

Description
===========

Fetch OpenLDAP schema for sudo(8):

```console
% fetch -o /tmp/sudo.schema https://raw.githubusercontent.com/sudo-project/sudo/refs/heads/main/docs/schema.OpenLDAP
```

Convert it to LDIF:

```console
% echo "include /tmp/sudo.schema" > /tmp/ldap_sudo.conf
% mkdir /tmp/ldap_sudo
% slaptest -F /tmp/ldap_sudo -f /tmp/ldap_sudo.conf
% slapcat -n0 -F /tmp/ldap_sudo -l /tmp/ldap_sudo_slapd.ldif
% awk '/^dn: cn=.*sudo,cn=schema,cn=config/,/^$/' /tmp/ldap_sudo_slapd.ldif \
    | sed \
        -e 's/^\(dn: cn=\).*\(sudo,cn=schema,cn=config\)/\1\2/' \
        -e 's/^\(cn: \).*\(sudo\)/\1\2/' \
    | awk '/^dn: cn=sudo,cn=schema,cn=config/,/^structuralObjectClass/{if(/^structuralObjectClass/) next; print}' \
    > /tmp/sudo.ldif
```

Import `/tmp/sudo.ldif` using ldapadd(1). Verify:

```console
% ldapisearch -Q -b cn=config -s sub '(cn=*sudo)' dn | grep '^dn:'
dn: cn={3}sudo,cn=schema,cn=config
```
