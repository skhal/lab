# NAME

**setup** - setup of `ldap.nuc.lab.net`

# DESCRIPTION

The instructions explain how to setup OpenLDAP server to host Unix users and
groups for Single Sign On in other jails.

## Create a jail

Create a VNET jail `ldap.nuc.lab.net` with IP `192.168.1.90/24`:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p2.2 zroot/jail/container/ldap 
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
<summary>Message from cyrus-sasl-2.1.28_5</summary>

```
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
```

</details>

<details>
<summary>Message from openldap26-client-2.6.10</summary>

```
The OpenLDAP client package has been successfully installed.

Edit
  /usr/local/etc/openldap/ldap.conf
to change the system-wide client defaults.

Try `man ldap.conf' and visit the OpenLDAP FAQ-O-Matic at
  http://www.OpenLDAP.org/faq/index.cgi?file=3
for more information.
```

</details>

<details>
<summary>Message from openldap26-server-2.6.10</summary>

```
The OpenLDAP server package has been successfully installed.

In order to run the LDAP server, you need to edit
  /usr/local/etc/openldap/slapd.conf
to suit your needs and add the following lines to /etc/rc.conf:
  slapd_enable="YES"
  slapd_flags='-h "ldapi://%2fvar%2frun%2fopenldap%2fldapi/ ldap://0.0.0.0/"'
  slapd_sockets="/var/run/openldap/ldapi"

Then start the server with
  /usr/local/etc/rc.d/slapd start
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
