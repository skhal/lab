# NAME

**setup** - setup of `ldap.nuc.lab.net`

# DESCRIPTION

Setup of `ldap.nuc.lab.net` describe how to:

  * Bootstrap OpenLDAP server slapd(8) ready for `dc=ldap,dc=net` database.
  * Configure ACL to give root user full access to the configuration `cn=config`
    database over local connections using Unix sockets.

> [!WARNING]
> The setup is minimal, without traffic encryption. SSL/TLS encryption is coming
> shortly.

## Create a VNET jail

Create a VNET jail `ldap.nuc.lab.net` with IP `192.168.1.90/24`
([doc](https://github.com/skhal/lab/blob/84821678384d2a7b4b6daa9b4e1266dd56cc9264/cluster/net.lab.nuc/doc/jail.md#vnet-jail)) by running the following commands from the jail hosting node
`nuc.lab.net`:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p2.3 zroot/jail/container/ldap 
# cat /etc/jail.conf.d/ldap.conf 
ldap {
  $id = "90";
  $ip = "192.168.1.${id}/24";
  $epair = "epair${id}";
  $bridge = "bridge0";

  # Virtual Network (VNET)
  vnet;
  vnet.interface = "${epair}b";
  devfs_ruleset = 5;

  exec.prestart += "/sbin/ifconfig ${epair} create";
  exec.prestart += "/sbin/ifconfig ${epair}a up descr jail:${name}";
  exec.prestart += "/sbin/ifconfig ${bridge} addm ${epair}a up";

  exec.start  = "/sbin/ifconfig ${epair}b ${ip} up";
  exec.start += "/bin/sh /etc/rc";

  exec.poststop += "/sbin/ifconfig ${bridge} deletem ${epair}a";
  exec.poststop += "/sbin/ifconfig ${epair}a destroy";
}
```

Start the jail and confirm it has Internet connection:

```console
# service jail start ldap
# jexec ldap ping -c 1 google.com
PING google.com (142.250.191.238): 56 data bytes
64 bytes from 142.250.191.238: icmp_seq=0 ttl=115 time=11.181 ms

--- google.com ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 11.181/11.181/11.181/0.000 ms
```

## Bootstrap OpenLDAP server

We'll start minimal configuration of LDAP server without databases.

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

Update OpenLDAP server slapd(8) configuration.

<details>
<summary>Changes</summary>

  * **Schemas**: include `cosine.ldif` and `nis.ldif` schemas for users and
    groups.

    https://github.com/skhal/lab/blob/6ae08fbaed51ae922e8b4c454f77e402d8872e3c/cluster/net.lab.nuc.ldap/doc/slapd.ldif.diff#L6-L8

  * **Database `dc=lab,dc=net`**: set administrator and password from
    `slappasswd`:

    https://github.com/skhal/lab/blob/6ae08fbaed51ae922e8b4c454f77e402d8872e3c/cluster/net.lab.nuc.ldap/doc/slapd.ldif.diff#L16-L24

    Store the database under `/var/db/openldap-data/lab.net`:

    https://github.com/skhal/lab/blob/6ae08fbaed51ae922e8b4c454f77e402d8872e3c/cluster/net.lab.nuc.ldap/doc/slapd.ldif.diff#L28-L29

    Speed up lookups with indices:

    https://github.com/skhal/lab/blob/6ae08fbaed51ae922e8b4c454f77e402d8872e3c/cluster/net.lab.nuc.ldap/doc/slapd.ldif.diff#L32-L33

    Access Control List (ACL) to restrict updates to user profiles to users and
    administrator - noone can read passwords:

    https://github.com/skhal/lab/blob/6ae08fbaed51ae922e8b4c454f77e402d8872e3c/cluster/net.lab.nuc.ldap/doc/slapd.ldif.diff#L34-L42
</details>

Push [spad.ldif.diff](./spad.ldif.diff) patch to `ldap.nuc.lab.net` and apply it
to the default slapd(8) configuration:

```console
# patch -R /usr/local/etc/openldap/slapd.ldif < ~/slapd.ldif.diff
Hmm...  Looks like a unified diff to me...
The text leading up to this was:
--------------------------
|--- /usr/local/etc/openldap/slapd.ldif.sample  2025-08-09 02:08:12.000000000 -0500
|+++ /usr/local/etc/openldap/slapd.ldif 2025-09-26 17:54:04.466570000 -0500
--------------------------
Patching file /usr/local/etc/openldap/slapd.ldif using Plan A...
Hunk #1 succeeded at 39.
Hunk #2 succeeded at 70.
Hunk #3 succeeded at 82.
done
```

slapd(8) works with a configuration directory in
`/usr/local/etc/openldap/slapd.d` instead of a configuration file. Create the
folder and import the configuration file into it:

```console
# mkdir /usr/local/etc/openldap/slapd.d
# /usr/local/sbin/slapadd -n0 -F /usr/local/etc/openldap/slapd.d/ -l /usr/local/etc/openldap/slapd.ldif
```

Create a folder to store `dc=lab,dc=net` database:

```console
# mkdir /var/db/openldap-data/lab.net
```

slapd(8) will run as `ldap:ldap` user. Fix permissions to the server
configuration and databases:

```console
# chown -R ldap:ldap /var/db/openldap-data /usr/local/etc/openldap/slapd.d
# chmod -R 700 /var/db/openldap-data /usr/local/etc/openldap/slapd.d
```

Configure and start slapd(8) service:

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
# service slapd start
```

> [!NOTE]
> The `slapd_sockets` flag forces the `slapd` rc-script to
> [fix](https://github.com/freebsd/freebsd-ports/blob/c2991243dbb2dfc9f932d1560af12061ed998cf2/net/openldap26-server/files/slapd.in#L153)
> the owner and permissions for `slapd:slapd` user.

Verity slapd(8) runs:

```console
# sockstat -4 | grep slapd
ldap     slapd      73735 8   tcp4   192.168.1.90:389      *:*
# ls -l /var/run/openldap/
total 6
srw-rw-rw-  1 ldap ldap   0 Sep 26 18:25 ldapi
-rw-r--r--  1 ldap ldap 105 Sep 26 18:25 slapd.args
-rw-r--r--  1 ldap ldap   6 Sep 26 18:25 slapd.pid
```

Items to check:

  * The service runs as `ldap` user.
  * slapd(8) is listening on the jail's IP address on the default port `:389`.
  * There must be a server socket and `slapd.pid` under `/var/run/openldap` with
    owner set to `ldap:ldap`.

## Root ACL for `cn=config`

slapd(8) prevents access to the configuration database `cn=config` by default:

```
olcAccess: to * by * none
```

We'll grant access to the root user when connected locally using Unix sockets
only, using `-Y EXTERNAL -H ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi` by
modifying the configuration database.

> [!CAUTION]
> Do not edit `/usr/local/etc/openldap/slapd.d` files manually. It may break
> consistency of the configuration directory. Change `.ldif` file and then
> import it into the directory configuration instead. See
> [doc](https://openldap.org/doc/admin26/slapdconf2.html).

The idea is to convert the configuration directory into an LDIF file, do the
update, and re-create the configuration directory while the server is down. We
can use this technique in the future to make further updates to the server.

Stop the service.

```console
# service slapd stop
```

Dump the directory service to an LDIF file:

```console
# slapcat -n0 -l /tmp/slapd.ldif
```

The root user connected over Unix sockets has
`gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth` identity. Grant
permissions to manage the configuration database:

```console
# slapcat -n0 | diff -u - /tmp/slapd.ldif 
--- -   2025-09-26 20:52:07.357710000 -0500
+++ /tmp/slapd.ldif     2025-09-26 20:50:22.724090000 -0500
@@ -574,7 +574,10 @@
 dn: olcDatabase={0}config,cn=config
 objectClass: olcDatabaseConfig
 olcDatabase: {0}config
-olcAccess: {0}to *  by * none
+olcAccess: {0}to *
+  by dn.exact=gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth manage
+  by * break
+olcAccess: {0}to * by * none
 olcAddContentAcl: TRUE
 olcLastMod: TRUE
 olcLastBind: FALSE
```

Re-create the configuration database and fix the permissions for `ldap:ldap`
user:

```console
# rm -rf /usr/local/etc/openldap/slapd.d/*
# /usr/local/sbin/slapadd -n0 -F /usr/local/etc/openldap/slapd.d/ -l /tmp/slapd.ldif
# chmod -R 700 /var/db/openldap-data /usr/local/etc/openldap/slapd.d
# chown -R ldap:ldap /var/db/openldap-data /usr/local/etc/openldap/slapd.d
```

Start the service:

```console
# service slapd start
```

Validate that the root user has access to the service:

> [!TIP]
> Use the following tcsh(1) alias to speed up search commands:<br/>
> `alias ldapisearch /usr/local/bin/ldapsearch -Y EXTERNAL -H ldapi://%2Fvar%2Frun%2Fopenldap%2Fldapi`

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
