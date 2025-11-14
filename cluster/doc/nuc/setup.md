NAME
====

**setup** - setup FreeBSD

Summary
=======

The setup instructions configure FreeBSD to host virtual nodes using classic, thick [jails](https://docs.freebsd.org/en/books/handbook/jails/).

The server has minimal setup:

-	Networking: two bridges for jails to separate internal and Internet traffic.
-	Security: packet filtering
-	Packages: doas, vim-tiny
-	Users: an operator user `op` to manage the system with SSH access.

The jails have the following setup:

-	Networking: full networking stack virtualisation through VNET.
-	Security: packet filtering
-	Packages: minimal set of packages to run jailed services or unrestricted for development jails.
-	Users: LDAP

Running jails:

-	Infrastructure:
	-	dns - BIND server
	-	ldap - OpenLDAP server
-	Development:
	-	dev - FreeBSD environment
	-	jammy - Linux environment via FreeBSD [Linux Binary Compatibility](https://docs.freebsd.org/en/books/handbook/linuxemu/)

Bootstrap
=========

Boot menu
---------

Disable boot menu to speed up restarts:

```console
# echo 'autoboot_delay="0"' >> /boot/loader.conf
```

System Update
-------------

Update to the latest patch available:

```console
# uname -a
# freebsd-update fetch
# freebsd-update install
```

Verify the updates applied as expected by checking the Kernel version:

```console
# uname -a
```

Packages
--------

Use latest packages available:

```console
# mkdir /usr/local/etc/pkg
# mkdir /usr/local/etc/pkg/repos
# cat /usr/local/etc/pkg/repos/FreeBSD.conf
FreeBSD: { url: "pkg+http://pkg.FreeBSD.org/${ABI}/latest" }
```

Update packages:

```console
# pkg update
# pkg install
```

Users
-----

Create a system operator:

```console
# pw groupadd \
    -g 1000 \
    -n op
# pw useradd \
    -c 'System Operator' \
    -d /home/op \
    -g op \
    -G wheel \
    -m \
    -n op \
    -s /bin/tcsh \
    -u 1000 \
    -w no
# passwd -l op
```

Change root-user shell to tcsh(1):

```console
# chsh -s /bin/tcsh root
```

Pluggable Authentication Modules (PAM)
--------------------------------------

Let PAM initialize user home folder.

```console
# pkg install pam_mkhomedir
```

Place the

```console
# grep pam_mkhomedir /etc/pam.d/login
session         required        /usr/local/lib/pam_mkhomedir.so
```

Resource configuration
----------------------

Ref: [rc](./rc.md)

Break monolith `/etc/rc.conf` into per-service configuration file to isolate service flags.

```console
# sysrc -s hostname -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/hostname /usr/local/etc/rc.conf.d/hostname
```

Move appropriate flags:

```console
# # add flag
# sysrc -f /etc/rc.conf.d/hostname hostname_foo=...
hostname_foo: ... -> ...
# # remove flag
# sysrc -x hostname_foo=...
```

Per-service configurations:

```console
# head /etc/rc.conf.d/*
==> /etc/rc.conf.d/hostname <==
hostname="nuc.lab.net"

==> /etc/rc.conf.d/moused <==
moused_nondefault_enable="no"

==> /etc/rc.conf.d/network <==
ifconfig_igc0="DHCP"

==> /etc/rc.conf.d/ntpd <==
ntpd_enable="yes"
ntpd_sync_on_start="yes"

==> /etc/rc.conf.d/powerd <==
powerd_enable="yes"

==> /etc/rc.conf.d/routing <==
defaultrouter="192.168.1.1"

==> /etc/rc.conf.d/sshd <==
sshd_enable="yes"

==> /etc/rc.conf.d/zfs <==
zfs_enable="yes"
```

Catch-all RC configuration:

```console
# cat /etc/rc.conf
clear_tmp_enable="YES"
dumpdev="AUTO"
```

Services
========

The instructions below skip basic service setup, covered in the Bootstrap section.

Network time protocol (NTP)
---------------------------

Limit open ports by ntpd:

```
% tail -n 3 /etc/ntp.conf
interface ignore wildcard
interface listen 127.0.0.1
interface listen igc0
```

Restart the service and verify open ports:

```console
% doas service ntpd restart
% doas sockstat -l4 | grep ntp
ntpd     ntpd       30027 20  udp4   192.168.1.102:123     *:*
ntpd     ntpd       30027 22  udp4   127.0.0.1:123         *:*
```

Packet filter (firewall)
------------------------

Use pf(4) to filter network traffic along with pflog(4) to access logs:

```console
% head /etc/rc.conf.d/pf*
==> /etc/rc.conf.d/pf <==
pf_enable="yes"

==> /etc/rc.conf.d/pflog <==
pflog_enable="yes"
```

pf(4) does not start without a configuration file even though it was enabled.

There are many helpful examples of basic pf(4) configuration available in `/usr/shared/examples/pf`. A minimal setup might be:

```
ext_if = "igc0"

tcp_services = "{ domain, http, https, ntp, ssh }"
udp_services = "{ domain, ntp }"

icmp_types = "{ echoreq, unreach }"

# deny all incoming traffic
block in all

# allow connections created by the system, retain the state.
pass out proto tcp to any port $tcp_services keep state
pass proto udp to any port $udp_services keep state

# debugging traffic
pass inet proto icmp all icmp-type $icmp_types keep state

# ssh
pass in inet proto tcp to $ext_if port ssh
```

Use pfctl(8) to manage PF:

```
pfctl -e                           # enable PF
pfctl -d                           # disable PF
pfctl -F all -f /etc/pf.conf       # flush NAT, filter, state, table rules, reload /etc/pf.conf
pfctl -s [ rules | nat | states ]  # report on filter rules, NAT rules, state tables
pfctl -si                          # report stats
pfctl -vn -f /etc/pf.conf          # check /etc/pf.conf for errors, do not load ruleset
```

Monitor traffic with pftop(1).

Routing
-------

The resolvconf(8) service manages `/etc/resolv.conf`, automatically created for DHCP networks. It picks up `deafultgateway` from the `routing` service.

Verity the configuration:

```console
% cat /etc/resolv.conf
# Generated by resolvconf
search localdomain
nameserver 192.168.1.1
```

Use drill(1) to test DNS forward and backward resolutions:

```console
% drill freebsd.org
;; ->>HEADER<<- opcode: QUERY, rcode: NOERROR, id: 24872
;; flags: qr rd ra ; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; QUESTION SECTION:
;; freebsd.org. IN  A

;; ANSWER SECTION:
freebsd.org.    3203INA96.47.72.84

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:

;; Query time: 2 msec
;; SERVER: 192.168.1.1
;; WHEN: Fri Nov 14 12:07:22 2025
;; MSG SIZE  rcvd: 45
% drill -x 96.47.72.84
;; ->>HEADER<<- opcode: QUERY, rcode: NOERROR, id: 12057
;; flags: qr rd ra ; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; QUESTION SECTION:
;; 84.72.47.96.in-addr.arpa.    IN  PTR

;; ANSWER SECTION:
84.72.47.96.in-addr.arpa.   3187    IN  PTR wfe0.nyi.freebsd.org.

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:

;; Query time: 2 msec
;; SERVER: 192.168.1.1
;; WHEN: Fri Nov 14 12:07:38 2025
;; MSG SIZE  rcvd: 76
```

Secure Shell
------------

Force SSH to listen on local port and restrict users that can SSH into the host:

```console
% cat /etc/rc.conf.d/sshd
sshd_enable="yes"
sshd_flags="-o ListenAddress=192.168.1.102 -o AllowUsers=op"
```

Restart the service and verify work:

```
% doas sockstat -l4 | grep sshd
root     sshd       39929 7   tcp4   192.168.1.102:22      *:*
```

Syslog
------

Close network socket:

```console
% cat /etc/rc.conf.d/syslogd
syslogd_flags="-ss"
```

Restart the service and verify work:

```
% doas sockstat -l4 | grep syslog
```

FreeBSD comes with logs rotation service newsyslog(8) enabled by default:

```console
% sysrc -A | grep newsyslog
newsyslog_enable: YES
newsyslog_flags: -CN
```

The default configuration `-CN` only forces newsyslog(8) to create missing log files. Use crontab(1) to rotate files every hour and use datetime instead of numeric suffix with `crontab -e`:

```console
# crontab -l
SHELL=/bin/sh
PATH=/sbin:/bin:/usr/sbin:/usr/bin

# Fields order
# minute hour mday month wday command

# Rotate log files every hour, if necessary
0 * * * * newsyslog -t DEFAULT
```

### Compression

newsyslog(8) uses bzip2(1) compression by default, newsyslog.conf(5). Notice "J" in the flags column of `/etc/newsyslog.conf`. This compression runs in the user land.

When setup with ZFS, FreeBSD configures ZFS to also use compression:

```console
% zfs get compression zroot/var/log
NAME           PROPERTY     VALUE           SOURCE
zroot/var/log  compression  lz4             inherited from zroot
```

In fact, this compression applies to all datasets under `zroot`. Check the compression efficiency after running the system for a while to let log files grow:

```console
% zfs get used,refcompressratio /var/log
NAME           PROPERTY          VALUE     SOURCE
zroot/var/log  used              292K      -
zroot/var/log  refcompressratio  3.61x     -
```

Given this setup, It is optional to disable compression all together in newsyslog(8) and let ZFS compress the logs including rotated files to save resources. However, copying the rotated file over the network may be inefficient unless one first compresses the file.

### Console logging

`rc(8)` scripts log messages using `logger(1)`, backed to `syslog(3)`. This writes messages to console and any other sinks configured in `systlogd(8)`. Unfortunately, `syslogd(8)` starts at some later stage in `rc(8)` sequence:

```console
% rcorder -p | cat -n | grep -C 1 syslogd
    29	/etc/rc.d/accounting /etc/rc.d/cleartmp /etc/rc.d/devfs /etc/rc.d/dmesg /etc/rc.d/gptboot /etc/rc.d/hostapd /etc/rc.d/mdconfig2 /etc/rc.d/motd /etc/rc.d/newsyslog /etc/rc.d/os-release /etc/rc.d/virecover /etc/rc.d/wpa_supplicant
    30	/etc/rc.d/syslogd
    31	/etc/rc.d/auditd /etc/rc.d/bsnmpd /etc/rc.d/hastd /etc/rc.d/ntpdate /etc/rc.d/power_profile /etc/rc.d/pwcheck /etc/rc.d/savecore /etc/rc.d/watchdogd
```

In result, rc-messages prior to `syslogd(8)` are only availalbe in console. The first ones are from `auditd` service in `/var/log/messages`:

```
Aug 29 14:40:08 nuc kernel: em0: link state changed to UP
Aug 29 14:40:08 nuc root[26487]: /etc/rc: DEBUG: checkyesno: auditd_enable is set to NO.
```

Configure `syslog(3)` to log all console messages:

```console
% grep console.info /etc/syslog.conf
console.info					/var/log/console.log
```

APPLICATIONS
============

Ref: https://docs.freebsd.org/en/books/handbook/ports/

Installed packages:

```console
% pkg prime-list
doas
gpu-firmware-intel-kmod-skylake
pkg
rsync
vim-tiny
wifi-firmware-iwlwifi-kmod-8000
```

Doas
----

Let the operator monitor the system:

```console
# pkg install doas
# cat <<eof >/usr/local/etc/doas.conf
permit nopass op cmd sockstat
eof
```

Vim
---

There are different pre-built configurations of [vim](https://www.freshports.org/editors/vim/). `vim-tiny` only includes binary without runtime files. It makes `vim-tiny` suitable for minimal installations.

```console
# pkg install vim-tiny
```
