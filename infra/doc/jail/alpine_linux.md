<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**alpine-linux** - enjail [Alpine Linux](https://www.alpinelinux.org)

# DESCRIPTION

**TL;DR** create a thin jail with Alpine Linux jail without FreeBSD jailed
environment.

```console
# zfs create zroot/usr/jail/alpine
# fetch -o /tmp https://dl-cdn.alpinelinux.org/alpine/latest-stable/releases/x86_64/alpine-minirootfs-3.23.3-x86_64.tar.gz
# tar -C /usr/jail/alpine -xzf /tmp/alpine-minirootfs-3.23.3-x86_64.tar.gz
```

Configure:

```console
# echo 'nameserver 10.0.0.2' >/usr/jail/alpine/etc/resolv.conf
# echo "auto lo" > /usr/jail/alpine/etc/network/interfaces
# chroot /usr/jail/alpine env PS1='alpine # ' /bin/sh
alpine # apk update
alpine # apk upgrade
alpine # apk add openrc
alpine # mkdir /run/openrc
alpine # touch /run/openrc/softlevel
```

Initialize passwords:

```console
# cd /usr/jail/alpine/etc/
# echo 'root::0:0::0:0:Charlie &:/root:/bin/sh' > master.passwd
# pwd_mkdb -d . -p master.passwd
```

Create a jail configuration:

```console
% cat /etc/jail.conf.d/alpine.conf
alpine {
  $id = "13";

  $gateway = "10.0.0.1";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.0.${id}/24";

  $bridges = "${bridge1}";
  $bridgeips = "${bridge1_ip}";

  host.hostname = "${name}.lab.net";
  path = "/usr/jail/${name}";

  depend  = "dns";
  depend += "ldap";

  vnet;

  # keep-sorted start
  allow.mount.devfs;
  allow.mount.fdescfs;
  allow.mount.linprocfs;
  allow.mount.linsysfs;
  allow.mount.procfs;
  allow.mount.tmpfs;
  allow.mount;
  allow.raw_sockets;
  # keep-sorted end

  mount += "devfs     $path/dev     devfs     rw 0 0";
  mount += "fdescfs   $path/dev/fd  fdescfs   rw,linrdlnk 0 0";
  mount += "linprocfs $path/proc    linprocfs rw 0 0";
  mount += "linsysfs  $path/sys     linsysfs  rw 0 0";
  mount += "tmpfs     $path/dev/shm tmpfs     rw,size=1g,mode=1777 0 0";
  mount.devfs;
  devfs_ruleset = 5;

  enforce_statfs = 1;

  exec.clean;
  exec.consolelog = "/var/log/jail_${name}.log";

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.created += "/sbin/route -j ${name} add default ${gateway}";
  exec.start    = "/sbin/openrc";

  persist;

  exec.stop     = "/sbin/openrc shutdown";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```

Start the jail and verify it works:

```console
% doas ifconfig -j alpine -ag epair
epair8b: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	description: jail:alpine:bridge1
	options=60000b<RXCSUM,TXCSUM,VLAN_MTU,RXCSUM_IPV6,TXCSUM_IPV6>
	ether 58:9c:fc:10:f6:4a
	inet 10.0.0.13 netmask 0xffffff00 broadcast 10.0.0.255
	groups: epair
	media: Ethernet 10Gbase-T (10Gbase-T <full-duplex>)
	status: active
	nd6 options=29<PERFORMNUD,IFDISABLED,AUTO_LINKLOCAL>
op@nuc:~ % doas route -j alpine get 0
   route to: default
destination: default
       mask: default
    gateway: 10.0.0.1
        fib: 0
  interface: epair8b
      flags: <UP,GATEWAY,DONE,STATIC>
 recvpipe  sendpipe  ssthresh  rtt,msec    mtu        weight    expire
       0         0         0         0      1500         1         0
% doas jexec alpine ping -c 1 alpinelinux.org
PING alpinelinux.org (213.219.36.190): 56 data bytes
64 bytes from 213.219.36.190: seq=0 ttl=44 time=108.967 ms

--- alpinelinux.org ping statistics ---
1 packets transmitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 108.967/108.967/108.967 ms
```

## Basics

```console
# apk add shadow shadow-doc tcsh tcsh-doc
# chsh -s /bin/tcsh root
# apk add mandoc mandoc-apropos less less-doc
# apk add patch patch-doc
# apk add drill
# apk add neovim
```

Many tools will have man(1) pages available in `foo-doc` package, e.g.:

```console
# apk add wget-doc
```

## Certs

Copy CA certificate to the jail:

```console
# cp /usr/jail/ssl/usr/local/share/certs/ca.crt /usr/jail/alpine/usr/local/share/ca-certificates/
# chroot /usr/jail/alpine env PS1='alpine # ' /bin/sh
alpine # sha256sum /etc/ssl/certs/ca-certificates.crt
766392c21c0baf5fa722cb309dc576b89d9fb3323dd32aa45a939dd575db6d1c  /etc/ssl/certs/ca-certificates.crt
alpine # update-ca-certificates
alpine # sha256sum /etc/ssl/certs/ca-certificates.crt
0e1052087075015f8e2ea0eab440de1a90f60c48cc76408d8c8e15f30263cca6  /etc/ssl/certs/ca-certificates.crt
```

## LDAP Client

```console
# apk add openldap-clients
# cd /tmp
# wget https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/ldap/client/openldap_ldap.conf.diff
# patch -lb -i /tmp/openldap_ldap.conf.diff /etc/openldap/ldap.conf
```

```console
# apk add nss-pam-ldapd nss-pam-ldapd-doc nss-pam-ldapd-openrc
# diff /etc/nslcd.conf{.orig,}
--- /etc/nslcd.conf.orig
+++ /etc/nslcd.conf
@@ -15,14 +15,14 @@
 #uri ldaps://127.0.0.1/
 #uri ldapi://%2fvar%2frun%2fldapi_sock/
 # Note: %2f encodes the '/' used as directory separator
-uri ldap://127.0.0.1/
+uri ldap://ldap.lab.net/

 # The LDAP version to use (defaults to 3
 # if supported by client library)
 #ldap_version 3

 # The distinguished name of the search base.
-base dc=example,dc=com
+base dc=lab,dc=net

 # The distinguished name to bind to the server with.
 # Optional: default is to bind anonymously.
@@ -59,12 +59,12 @@
 #idle_timelimit 3600

 # Use StartTLS without verifying the server certificate.
-#ssl start_tls
-#tls_reqcert never
+ssl start_tls
+tls_reqcert allow

 # CA certificates for server certificate verification
 #tls_cacertdir /etc/ssl/certs
-#tls_cacertfile /etc/ssl/ca.cert
+tls_cacertfile /etc/ssl/certs/ca-certificates.crt

 # Seed the PRNG if /dev/urandom is not provided
 #tls_randfile /var/run/egd-pool
# rc-update add nslcd
```

Add user folder mounts:

```console
# mkdir /home/op /home/skhalatyan
```

Add mount points for user home folders in Alpine jail configuration.

```console
# apk add openssh-server-pam
# addgroup -S sshd
# adduser -H -S -s /sbin/nologin -G sshd sshd
# diff /etc/ssh/sshd_config{.orig,}
--- /etc/ssh/sshd_config.orig
+++ /etc/ssh/sshd_config
@@ -85,7 +85,7 @@
 # If you just want the PAM account and session checks to run without
 # PAM authentication, then enable this but set PasswordAuthentication
 # and KbdInteractiveAuthentication to 'no'.
-#UsePAM no
+UsePAM yes

 #AllowAgentForwarding yes
 # Feel free to re-enable these if your use case requires them.
# diff /etc/init.d/sshd{.orig,}
--- /etc/init.d/sshd.orig
+++ /etc/init.d/sshd
@@ -12,7 +12,7 @@
 : "${cfgfile:=${SSHD_CONFIG:-"${SSHD_CONFDIR:-"/etc/ssh"}/sshd_config"}}"

 pidfile="${SSHD_PIDFILE:-"/run/$RC_SVCNAME.pid"}"
-command="${SSHD_BINARY:-"/usr/sbin/sshd"}"
+command="${SSHD_BINARY:-"/usr/sbin/sshd.pam"}"
 command_args="${command_args:-${SSHD_OPTS:-}}"

 required_files="$cfgfile"
```

Start the service with:

```console
# rc-service sshd start
# rc-service sshd status
 * status: started
```

It is hard to debug because none of the Linux networking tools work in this
setup - they can't access network stack:

```console
# busybox netstat -tuln
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State
netstat: /proc/net/tcp: No such file or directory
netstat: /proc/net/tcp6: No such file or directory
netstat: /proc/net/udp: No such file or directory
netstat: /proc/net/udp6: No such file or directory
# apk add iproute2
# ss -tln
Cannot open netlink socket: Protocol not supported
State     Recv-Q     Send-Q         Local Address:Port         Peer Address:Port
```

None of the tools generate logs except apk:

```console
# ls /var/log
apk.log
```

## Install Bazel

```console
# apk update
# apk upgrade
# echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
# apk update
# apk add bazel8
```
