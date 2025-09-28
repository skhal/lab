# NAME

**setup** - setup instructions for `snk.nuc.lab.net`

# DESCRIPTION

The setup runs a VNET jail `snk.nuc.lab.net` with Internet connection. It
uses LDAP for Single Sign On and SSH.

## Bootstrap VNET jail

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p2.3 zroot/jail/container/snk
# cat /etc/jail.conf.d/snk.conf
snk {
  $id = "112";
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

  depend = ldap;
}
# service jail start snk
```

## OpenLDAP client

```console
# pkg install openldap26-client
```

<details>
<summary>Messages from installed packages</summary>

```
=====
Message from cyrus-sasl-2.1.28_5:

--
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
=====
Message from openldap26-client-2.6.10:

--
The OpenLDAP client package has been successfully installed.

Edit
  /usr/local/etc/openldap/ldap.conf
to change the system-wide client defaults.

Try `man ldap.conf' and visit the OpenLDAP FAQ-O-Matic at
  http://www.OpenLDAP
```

</details>

Configure LDAP server default settings:

```console
# diff -u /usr/local/etc/openldap/ldap.conf{.sample,}
--- /usr/local/etc/openldap/ldap.conf.sample	2025-08-08 20:14:24.000000000 -0500
+++ /usr/local/etc/openldap/ldap.conf	2025-09-27 19:23:19.839450000 -0500
@@ -5,8 +5,10 @@
 # See ldap.conf(5) for details
 # This file should be world readable but not world writable.
 
-#BASE	dc=example,dc=com
-#URI	ldap://ldap.example.com ldap://ldap-provider.example.com:666
+BASE	dc=lab,dc=net
+URI	ldap://192.168.1.90
 
 #SIZELIMIT	12
 #TIMELIMIT	15
```

Verify it works:

```console
# ldapsearch -x dn | grep '^dn:'
dn: dc=lab,dc=net
dn: ou=people,dc=lab,dc=net
dn: ou=groups,dc=lab,dc=net
dn: uid=op,ou=people,dc=lab,dc=net
dn: cn=op,ou=groups,dc=lab,dc=net
```

## Single Sign On

We'll use Pluggable Authentication Modules (PAM) for a Single Sign On with
LDAP users and gruops.

```console
# pkg install pam_ldap nss_ldap
```

<details>
<summary>Messages from installed packages</summary>

```
=====
Message from nss_ldap-1.265_15:

--
The nss_ldap module expects to find its configuration files at the
following paths:

LDAP configuration:     /usr/local/etc/nss_ldap.conf
LDAP secret (optional): /usr/local/etc/nss_ldap.secret
=====
Message from pam_ldap-186_2:

--
Edit /usr/local/etc/ldap.conf in order to use this module.  Then
create a /usr/local/etc/pam.d/ldap with a line similar to the following:

login	auth	sufficient	/usr/local/lib/pam_ldap.so
```

</details>

Configure NSS with LDAP:

```console
# diff -u /usr/local/etc/nss_ldap.conf{.sample,}
--- /usr/local/etc/nss_ldap.conf.sample	2025-08-10 10:17:40.000000000 -0500
+++ /usr/local/etc/nss_ldap.conf	2025-09-27 19:45:24.726695000 -0500
@@ -12,10 +12,10 @@
 # space. How long nss_ldap takes to failover depends on
 # whether your LDAP client library supports configurable
 # network or connect timeouts (see bind_timelimit).
-host 127.0.0.1
+host 192.168.1.90
 
 # The distinguished name of the search base.
-base dc=padl,dc=com
+base dc=lab,dc=net
 
 # Another way to specify your LDAP server is to provide an
 # uri with the server name. This allows to use
@@ -87,7 +87,7 @@
 #pam_filter objectclass=account
 
 # The user ID attribute (defaults to uid)
-#pam_login_attribute uid
+pam_login_attribute uid
 
 # Search the root DSE for the password policy (works
 # with Netscape Directory Server)
```

Now PAM LDAP searches for entries with matching uid under LDAP base. It will
bind with the found record only if a single record is found.

Set Name Service Switch (NSS) to use local files and LDAP for user password
and group lookups:

```console
# diff -u /etc/nsswitch.conf{.orig,}
--- /etc/nsswitch.conf.orig	2025-09-27 19:37:20.601699000 -0500
+++ /etc/nsswitch.conf	2025-09-27 19:37:42.649046000 -0500
@@ -1,12 +1,12 @@
 #
 # nsswitch.conf(5) - name service switch configuration file
 #
-group: compat
+group: files ldap
 group_compat: nis
 hosts: files dns
 netgroup: compat
 networks: files
-passwd: compat
+passwd: files ldap
 passwd_compat: nis
 shells: files
 services: compat
```

More PAM configurations are under `/etc/pam.d/`.

## SSH

```console
# sysrc -f /etc/rc.conf.d/sshd sshd_flags="-o ListenAddress=192.168.1.112"
sshd_flags:  -> -o ListenAddress=192.168.1.112
# sysrc -f /etc/rc.conf.d/sshd sshd_enable=yes
sshd_enable: NO -> yes
```

Start the service

```console
# service sshd start
```

Verify it is running:

```console
# sockstat -4 | grep sshd
root     sshd        2662 7   tcp4   192.168.1.112:22      *:*
```
