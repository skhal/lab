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

Configure LDAP server default settings using
[`ldap.conf.diff`](./ldap.conf.diff) (push the file to the remote server):

https://github.com/skhal/lab/blob/0d47a18e2a4a68a668ff4164971160e17baba8bd/cluster/net.lab.nuc.snk/doc/ldap.conf.diff#L7-L10

```console
# patch /usr/local/etc/openldap/ldap.conf ~/ldap.conf.diff
Hmm...  Looks like a unified diff to me...
The text leading up to this was:
--------------------------
|--- /usr/local/etc/openldap/ldap.conf.sample 2025-08-08 20:14:24.000000000 -0500
|+++ /usr/local/etc/openldap/ldap.conf  2025-09-28 08:45:32.133654000 -0500
--------------------------
Patching file /usr/local/etc/openldap/ldap.conf using Plan A...
Hunk #1 succeeded at 5.
done
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

We'll use Name Switch Service (NSS) integration with LDAP to pull users and
passwords from LDAP after trying local users first.

```console
# pkg install nss_ldap
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
```

</details>

Use [`nss_ldap.conf.diff`](./nss_ldap.conf.diff) to patch nss_ldap(5)
configuration to point NSS to LDAP instance and define base for lookups.

https://github.com/skhal/lab/blob/0d47a18e2a4a68a668ff4164971160e17baba8bd/cluster/net.lab.nuc.snk/doc/nss_ldap.conf.diff#L7-L20

```console
# patch /usr/local/etc/nss_ldap.conf ~/nss_ldap.conf.diff
Hmm...  Looks like a unified diff to me...
The text leading up to this was:
--------------------------
|--- /usr/local/etc/nss_ldap.conf.sample  2025-08-10 10:17:40.000000000 -0500
|+++ /usr/local/etc/nss_ldap.conf 2025-09-28 08:56:50.439652000 -0500
--------------------------
Patching file /usr/local/etc/nss_ldap.conf using Plan A...
Hunk #1 succeeded at 12.
Hunk #2 succeeded at 24.
done
```

Finally, let NSS use local files and LDAP to lookup users and groups in that
order by patching `nsswitch.conf` with
[`nsswitch.conf.diff`](./nsswitch.conf.diff):

https://github.com/skhal/lab/blob/0d47a18e2a4a68a668ff4164971160e17baba8bd/cluster/net.lab.nuc.snk/doc/nsswitch.conf.diff#L7-L14

```console
# patch /etc/nsswitch.conf ~/nsswitch.conf.diff 
Hmm...  Looks like a unified diff to me...
The text leading up to this was:
--------------------------
|--- /etc/nsswitch.conf.sample  2025-09-28 09:01:28.828836000 -0500
|+++ /etc/nsswitch.conf 2025-09-28 09:02:30.649736000 -0500
--------------------------
Patching file /etc/nsswitch.conf using Plan A...
Hunk #1 succeeded at 1.
done
```

Now PAM LDAP searches for entries with matching uid under LDAP base. It will
bind with the found record only if a single record is found.

```console
# getent group op
op:*:1000:op
# getent passwd op
op:*:1000:1000:Operator:/home/op:/bin/tcsh
```

## SSH

Use Pluggable Authentication Modules to configure SSH with LDAP:

```console
# pkg install pam_ldap
```

<details>
<summary>Messages from installed packages</summary>

```
=====
Message from pam_ldap-186_2:

--
Edit /usr/local/etc/ldap.conf in order to use this module.  Then
create a /usr/local/etc/pam.d/ldap with a line similar to the following:

login	auth	sufficient	/usr/local/lib/pam_ldap.so
```

</details>

> [!IMPORTANT]
> PAM LDAP uses a different `ldap.conf` file, located at
> `/usr/local/etc/ldap.conf`. Recall that LDAP client configuration used
> `/usr/local/etc/opendlap/ldap.conf` to run `ldapsearch` and other LDAP
> commands.

Patch LDAP configuration for PAM:

```console
# patch /usr/local/etc/ldap.conf ~/pam_ldap.conf.diff 
Hmm...  Looks like a unified diff to me...
The text leading up to this was:
--------------------------
|--- /usr/local/etc/ldap.conf.sample  2025-08-09 11:23:34.000000000 -0500
|+++ /usr/local/etc/ldap.conf 2025-09-28 09:25:14.265522000 -0500
--------------------------
Patching file /usr/local/etc/ldap.conf using Plan A...
Hunk #1 succeeded at 12.
Hunk #2 succeeded at 24.
done
```

PAM policies use pam.conf(5) format. Per-service policies are under
`/etc/pam.d/` or `/usr/local/etc/pam.d/` for installed packages respectively.

Use [`pam_sshd.diff`](./pam_sshd.diff) to patch PAM policy for SSH server:

```console
# patch /etc/pam.d/sshd ~/pam_sshd.diff 
Hmm...  Looks like a unified diff to me...
The text leading up to this was:
--------------------------
|--- /etc/pam.d/sshd.sample 2025-09-28 09:15:43.174489000 -0500
|+++ /etc/pam.d/sshd  2025-09-28 09:36:34.631403000 -0500
--------------------------
Patching file /etc/pam.d/sshd using Plan A...
Hunk #1 succeeded at 6.
done
```

Configure SSH server:

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

Verify that SSH server runs and it works using an `op` user from LDAP:

```console
# sockstat -4 | grep sshd
root     sshd        2662 7   tcp4   192.168.1.112:22      *:*
# ssh op@192.168.1.112 whoami
(op@192.168.1.112) Password: 
Could not chdir to home directory /home/op: No such file or directory
op
```
