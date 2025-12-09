NAME
====

**ldap server setup** - how to setup an LDAP server

DESCRIPTION
===========

Install
-------

Install OpenLDAP server:

```console
# pkg install openldap26-server
```

<details>
<summary>Messages from installed packages</summary>

```
===
Message from cyrus-sasl-2.1.28_5

You can use sasldb2 for authentication, to add users use:

  saslpasswd2 -c username

If you want to enable SMTP AUTH with the system Sendmail, read
Sendmail.README

NOTE: This port has been compiled with a default pwcheck_method of
      auxprop.  If you want to authenticate your user by /etc/passwd,
      PAM or LDAP, install ports/security/cyrus-sasl2-saslauthd and
      set sasl_pwcheck_method to saslauthd after installing the
      Cyrus-IMAPd 2.X port.  You should also check the
      /usr/local/lib/sasl2/*.conf files for the correct
      pwcheck_method.
      If you want to use GSSAPI mechanism, install
      ports/security/cyrus-sasl2-gssapi.
      If you want to use SRP mechanism, install
      ports/security/cyrus-sasl2-srp.
      If you want to use LDAP auxprop plugin, install
      ports/security/cyrus-sasl2-ldapdb.

===
Message from openldap26-client-2.6.10

The OpenLDAP client package has been successfully installed.

Edit /usr/local/etc/openldap/ldap.conf
to change the system-wide client defaults.

Try `man ldap.conf' and visit the OpenLDAP FAQ-O-Matic at
http://www.OpenLDAP.org/faq/index.cgi?file=3
for more information.

===
Message from openldap26-server-2.6.10

The OpenLDAP server package has been successfully installed.

In order to run the LDAP server, you need to edit
`/usr/local/etc/openldap/slapd.conf`
to suit your needs and add the following lines to /etc/rc.conf:

  slapd_enable="YES"
  slapd_flags='-h "ldapi://%2fvar%2frun%2fopenldap%2fldapi/ ldap://0.0.0.0/"'
  slapd_sockets="/var/run/openldap/ldapi"

Then start the server with
  usr/local/etc/rc.d/slapd start
or reboot.

Try `man slapd' and the online manual at
http://www.OpenLDAP.org/doc/
for more information.

slapd runs under a non-privileged user id (by default `ldap'),
see /usr/local/etc/rc.d/slapd for more information.

PLEASE NOTE:

Upgrading from openldap26-server 2.4 to 2.5 requires a full dump
and reimport of database.

Starting from openldap26-server 2.4.59_3, automatic data dumps
are saved at /var/backups/openldap when shutting down slapd.

Please refer to OpenLDAP Software 2.5 Administrator's Guide at
  https://www.openldap.org/doc/admin25/appendix-upgrading.html
for additional upgrade instructions.
```
</details>

Configure
---------

A running OpenLDAP server slapd(8) holds the configuration in a configuration database under `/usr/local/etc/openldap/slapd.d`. The installation package comes with a default configuration file `/usr/local/etc/openldap/slapd.ldif`, which is to be imported into the configuration database with `slapadd(8)` tool.

First, modify the default configuration settings in `slapd.ldif` configuration file using [`slapd.ldif.diff`](./slapd.ldif.diff). It makes the following changes:

-	Load additional LDAP schemas to manage users & groups: `cosine.ldif` and`nis.ldif`.
-	Create a database `dc=lab,dc=net`:
	-	Configure the administrator's account `olcRootDN` and password`olcRootPW`. Use slappasswd(8) to generate encoded password.
	-	Isolate the database under `/var/db/opendal-data/lab.net`.
	-	Add indices to speed up user account lookups.
	-	Restrict DB access using Access Control Lists (ACLs)

Patch the configuration file ([server/slapd.ldif.diff](server/slapd.ldif.diff)\):

```console
# fetch -o /tmp https://github.com/skhal/lab/raw/refs/heads/main/cluster/doc/ldap/server/slapd.ldif.diff
# patch -R /usr/local/etc/openldap/slapd.ldif < /tmp/slapd.ldif.diff
```

The configuration file is ready, import it into a configuration database:

```console
# mkdir /usr/local/etc/openldap/slapd.d
# mkdir /var/db/openldap-data/lab.net
# /usr/local/sbin/slapadd -n0 -F /usr/local/etc/openldap/slapd.d/ -l /usr/local/etc/openldap/slapd.ldif
```

Bootstrap filesystem
--------------------

slapd(8) will run as `ldap:ldap` user. Fix permissions to the server configuration and databases:

```console
# chown -R ldap:ldap /var/db/openldap-data /usr/local/etc/openldap/slapd.d
# chmod -R 700 /var/db/openldap-data /usr/local/etc/openldap/slapd.d
```

Service
-------

Configure slapd(8) service:

```console
# mkdir /usr/local/etc/rc.conf.d
# sysrc -f /usr/local/etc/rc.conf.d/slapd slapd_enable=yes
slapd_enable:  -> yes
# sysrc -f /usr/local/etc/rc.conf.d/slapd slapd_sockets="/var/run/openldap/ldapi"
slapd_sockets:  -> /var/run/openldap/ldapi
# sysrc -f /usr/local/etc/rc.conf.d/slapd slapd_flags="-h 'ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi ldap://192.168.1.90'"
slapd_flags:  -> -h 'ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi ldap://192.168.1.90:389/'
# sysrc -f /usr/local/etc/rc.conf.d/slapd slapd_cn_config=yes
slapd_cn_config:  -> yes
```

> [!NOTE] The `slapd_sockets` flag forces the `slapd` rc-script to [fix](https://github.com/freebsd/freebsd-ports/blob/c2991243dbb2dfc9f932d1560af12061ed998cf2/net/openldap26-server/files/slapd.in#L153) the owner and permissions for `slapd:slapd` user.

Start the service

```console
# service slapd start
```

Verify that:

-	the service runs under `ldap` user.
-	slapd(8) listens on a single IP address at port `:389`.
-	A server socket and other files under `/var/run/openldap` have `ldap:ldap` ownership.

	```console
	# sockstat -l4
	USER     COMMAND    PID   FD  PROTO  LOCAL ADDRESS         FOREIGN ADDRESS
	ldap     slapd      56701 7   tcp4   10.0.1.3:389          *:*
	# ls -l /var/run/openldap/
	total 6
	srw-rw-rw-  1 ldap ldap   0 Sep 26 18:25 ldapi
	-rw-r--r--  1 ldap ldap 105 Sep 26 18:25 slapd.args
	-rw-r--r--  1 ldap ldap   6 Sep 26 18:25 slapd.pid
	```

Root ACL for `cn=config`
------------------------

By default, slapd(8) prevents access to `cn=config`, where all the server configurations reside:

```
olcAccess: to * by * none
```

We want to grant permissions to the root user when connected to slapd(8) locally via Unix sockets, e.g.:

```console
# ldapsearch -H ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi -Y EXTERNAL ...
```

The root user has the following identity:

```
gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth
```

The process consists of the following steps while the service is down:

-	Dump the configuration database into a temporary LDIF configuration file.
-	Update the configuration file.
-	Re-create the configuration database.

> [!CAUTION] Do not edit `/usr/local/etc/openldap/slapd.d` files manually. It may break consistency of the configuration directory. Change `.ldif` file and then import it into the directory configuration instead. See [doc](https://openldap.org/doc/admin26/slapdconf2.html).

Stop the service:

```console
# service slapd stop
```

Dump the directory service, including the configuration, to a temporary file:

```console
# slapcat -n0 -l /tmp/slapd.ldif
```

Update access to the configuration database:

```console
# slapcat -n0 | diff -u - /tmp/slapd.ldif
--- -   2025-09-26 20:52:07.357710000 -0500
+++ /tmp/slapd.ldif     2025-09-26 20:50:22.724090000 -0500
@@ -574,7 +574,8 @@
 dn: olcDatabase={0}config,cn=config
 objectClass: olcDatabaseConfig
 olcDatabase: {0}config
-olcAccess: {0}to *  by * none
+olcAccess: {0}to * by dn.exact=gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth manage by * break
+olcAccess: {0}to * by * none
 olcAddContentAcl: TRUE
 olcLastMod: TRUE
 olcLastBind: FALSE
```

Re-create the configuration database:

```console
# rm -rf /usr/local/etc/openldap/slapd.d/*
# /usr/local/sbin/slapadd -n0 -F /usr/local/etc/openldap/slapd.d/ -l /tmp/slapd.ldif
```

Fix the configuration directory permissions: directory permissions: directory permissions: directory permissions:

```console
# chmod -R 700 /var/db/openldap-data /usr/local/etc/openldap/slapd.d
# chown -R ldap:ldap /var/db/openldap-data /usr/local/etc/openldap/slapd.d
```

Start the service:

```console
# service slapd start
```

Verify:

> [!TIP] Use the following tcsh(1) alias to speed up search commands:<br/> `alias ldapisearch /usr/local/bin/ldapsearch -Y EXTERNAL -H ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi`

```console
# ldapisearch -b cn=config dn | grep '^dn:'
SASL/EXTERNAL authentication started
SASL username: gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth
SASL SSF: 0
dn: cn=config
dn: cn=module{0},cn=config
dn: cn=schema,cn=config
dn: cn={0}core,cn=schema,cn=config
dn: cn={1}cosine,cn=schema,cn=config
dn: cn={2}nis,cn=schema,cn=config
dn: olcDatabase={-1}frontend,cn=config
dn: olcDatabase={0}config,cn=config
dn: olcDatabase={1}mdb,cn=config
dn: olcDatabase={2}monitor,cn=config
```

Management
==========

Any client can update LDAP server contents with instructions in LDIF format using ldapadd(1) and ldapmodify(1) commands that pick up LDAP URI from the configuration (use `-H ldap://ldap.lab.net` flag to explicitly set URI):

```console
% ldapadd -x -W -D 'cn=op,dc=lab,dc=net' -f /tmp/foo.ldif
% ldapmodify -x -W -D 'cn=op,dc=lab,dc=net' -f /tmp/foo.ldif
```

Change user password with ldappasswd(1):

```console
% ldappasswd -x -W -D 'cn=op,dc=lab,dc=net' -S 'uid=foo,ou=people,dc=lab,dc=net'
```
