## Sudo schema

Source: https://raw.githubusercontent.com/sudo-project/sudo/refs/heads/main/docs/schema.OpenLDAP

```console
% fetch -o ./openldap/schema/sudo.schema https://raw.githubusercontent.com/sudo-project/sudo/refs/heads/main/docs/schema.OpenLDAP
openldap/schema/sudo.schema                           2499  B   14 MBps    00s
```

Convert `.schema` to `.ldif`:

```console
% setenv SUDO_SCHEMA /path/to/sudo.schema
% setenv SUDO_LDIF ${SUDO_SCHEMA:s/.schema/.ldif/} 
% mkdir /tmp/ldap_sudo
% echo "include ${SUDO_SCHEMA}" > /tmp/ldap_sudo.conf
% slaptest -F /tmp/ldap_sudo -f /tmp/ldap_sudo.conf
% slapcat -n0 -F /tmp/ldap_sudo -l /tmp/ldap_sudo_slapd.ldif
% awk '/^dn: cn=.*sudo,cn=schema,cn=config/,/^$/' /tmp/ldap_sudo_slapd.ldif \
    | sed \
        -e 's/^\(dn: cn=\).*\(sudo,cn=schema,cn=config\)/\1\2/' \
        -e 's/^\(cn: \).*\(sudo\)/\1\2/' \
    | awk '/^dn: cn=sudo,cn=schema,cn=config/,/^structuralObjectClass/{if(/^structuralObjectClass/) next; print}' \
    > ${SUDO_LDIF}
```

Load schema into the server and confirm:

```console
% ldapadd -H ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi -Y EXTERNAL -f ${SUDO_LDIF}
% ldapisearch -Q -b cn=config -s sub '(cn=*sudo)' dn | grep '^dn:'
dn: cn={3}sudo,cn=schema,cn=config
```
