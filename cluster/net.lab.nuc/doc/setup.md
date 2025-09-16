# NAME

**setup** - basic setup of `nuc.lab.net`


# DESCRIPTION

Basic setup brings a fresh installation of FreeBSD to minimal setup,
suitable for operations. It covers:

  * Create a user `op` to operate the OS through `wheel` group membership.
  * Move service configs from global `/etc/rc.conf` to per-service configuration
    file under `/etc/rc.conf.d/`.
  * Install few applications to manage the system.


# USERS

Create an operator user `op` to manage the node:

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

Isolate user home folder in a dedicated ZFS dataaset and let the user
manage it:

```console
# zfs create zroot/home/op
# zfs allow -u op create,destroy,mount,snapshot zroot/home/op
# sh -c 'for f in /usr/share/skel/*; do cp $f ~/.${f##*/dot.}; done'
# chown -R op:op /home/op
```

Change root-user shell to tcsh(1):

```console
# chsh -s /bin/tcsh root
```


# SERVICES

FreeBSD boot process uses `init(8)`. It triggers `rc(8)` to start services.
Every service has a script in one of `rc.d/` folders: standard location
`/etc/rc.d` and folders set by `local_startup` flag:

```console
% sysrc local_startup
local_startup: /usr/local/etc/rc.d
```

Service scripts define dependencies, rc-variables, actions (start, stop, etc.)
and function to execute the actions.

Use `rcorder(8)` to dump services dependency graph in topological order (`-p`
groups services that can start in parallel`):

```console
% rcorder -p /etc/rc.d/* | head -n 2
/etc/rc.d/dhclient /etc/rc.d/dumpon /etc/rc.d/dnctl /etc/rc.d/natd /etc/rc.d/sysctl
/etc/rc.d/ddb /etc/rc.d/hostid
```

Service scripts use variables, aka flags, set in `rc.conf(5)` files, loaded
in the following order (the last loaded value wins):

  * A default value is set in `/etc/default/rc.conf` with optional override
    from `/etc/default/vendor.conf` (if exists).
  * `/etc/rc.conf` and `/etc/rc.conf.local` (legacy) are global containers,
    loaded by all services including `rc(8)` itself.
  * `<dir>/rc.conf.d/<name>` is only loaded by the service `<name>`, where
    `<dir>/` is either the standard location `/etc` or folders listed in
    `local_startup` with `rc.d/` suffix removed.

Use `sysrc(8)` to list supported configuration files for a given service:

```console
% sysrc -s hostname -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/hostname /usr/local/etc/rc.conf.d/hostname
```

> [!WARNING]
> `/etc/rc.conf` is a global container of flags, shared between all services
> and `rc(8)`. In order to limit the visibility of flags, place flags into
> shared files under `<dir>/rc.conf.d/`.
>
> For example, `dhclient` and `netif` services share DHCP settings via
> `/etc/rc.conf.d/network`.
>
> ```console
> % sysrc -s dhclient -l
> /etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/dhclient /usr/local/etc/rc.conf.d/dhclient /etc/rc.conf.d/network /usr/local/etc/rc.conf.d/network
> % sysrc -s netif -l
> /etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/network /usr/local/etc/rc.conf.d/network /etc/rc.conf.d/netif /usr/local/etc/rc.conf.d/netif
> ```

The instructions below move flags from `/etc/rc.conf` to per-service
configuration file under `/etc/rc.conf.d/`. It preserves `/usr/local/etc` 
for services installed by ports and packages.

## Hostname

```console
% sysrc -s hostname -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/hostname /usr/local/etc/rc.conf.d/hostname
# sysrc -x hostname
# sysrc -f /etc/rc.conf.d/hostname hostname="nuc.lab.net"
# service hostname restart
```

## Moused

```console
% sysrc -s moused -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/moused /usr/local/etc/rc.conf.d/moused
# sysrc -x moused_nondefault_enable
# sysrc -f /etc/rc.conf.d/moused moused_nondefault_enable=no
moused_nondefault_enable: YES -> no
```

## Network

`netif` manages network interfaces. It shares some of the DHCP configurations
with `dhclient` via `/etc/rc.conf.d/network`, therefore keep configuration in
this file.

```console
% sysrc -s netif -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/network /usr/local/etc/rc.conf.d/network /etc/rc.conf.d/netif /usr/local/etc/rc.conf.d/netif
# sysrc -x ifconfig_em0
# sysrc -f /etc/rc.conf.d/network ifconfig_em0=DHCP
# shutdown -r now
```

It is best to reboot the host for all the services to pick up the network file.

### Wireless

Ref: https://docs.freebsd.org/en/books/handbook/network/#network-wireless

First, create a `wpa_supplicant.conf(5)` with wireless Service Set Identifier
(SSID) and Pre-Shared Key (PSK) (the contents are redacted):

```console
% cat /etc/wpa_supplicant.conf
network={
  ssid="<wifi-id>"
  bssid=<access-point-mac-address>
  psk="<password>"
}
```

Identify wireless network device interface:

```console
% sysctl net.wlan.devices
net.wlan.devices: iwm0
```

Create a wireless LAN on this interface:

```console
# sysrc -f /etc/rc.conf.d/network wlans_iwm0="wlan0"
wlans_iwm0:  -> wlan0
# sysrc -f /etc/rc.conf.d/network ifconfig_wlan0="WPA DHCP"
ifconfig_wlan0:  -> WPA DHCP
```

## Network Time server

```console
% sysrc -s ntpd -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/ntpd /usr/local/etc/rc.conf.d/ntpd
# sysrc -x ntpd_sync_on_start
# sysrc -x ntpd_enable
# sysrc -f /etc/rc.conf.d/ntpd ntpd_enable=yes
# sysrc -f /etc/rc.conf.d/ntpd ntpd_sync_on_start=yes
# service ntpd restart
```

## Routing

```console
% sysrc -s routing -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/routing /usr/local/etc/rc.conf.d/routing
# sysrc -f /etc/rc.conf.d/routing defaultrouter=192.168.1.1
# service routing restart
```

Check that resolver configuration is applied systemwide:

```console
% cat /etc/resolv.conf
# Generated by resolvconf
search localdomain
nameserver 192.168.1.1
```

Verify that DNS resolution uses router:

```console
 % drill freebsd.org
;; ->>HEADER<<- opcode: QUERY, rcode: NOERROR, id: 5988
;; flags: qr rd ra ; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0 
;; QUESTION SECTION:
;; freebsd.org.	IN	A

;; ANSWER SECTION:
freebsd.org.	3600	IN	A	96.47.72.84

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:

;; Query time: 25 msec
;; SERVER: 192.168.1.1
;; WHEN: Mon Sep  1 09:
```

## RC

Enable rc debug and info logging:

```console
# sysrc rc_debug=yes
# sysrc rc_info=yes
```

## Syslog

```console
% sysrc -s syslogd -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/syslogd /usr/local/etc/rc.conf.d/syslogd
# sysrc -x syslogd_flags
# sysrc -f /etc/rc.conf.d/syslogd syslogd_flags="-ss"
# service syslogd restart
```

`rc(8)` scripts log messages using `logger(1)`, backed to `syslog(3)`. This
writes messages to console and any other sinks configured in `systlogd(8)`.
Unfortunately, `syslogd(8)` starts at some later stage in `rc(8)` sequence:

```console
% rcorder -p | cat -n | grep -C 1 syslogd
    29	/etc/rc.d/accounting /etc/rc.d/cleartmp /etc/rc.d/devfs /etc/rc.d/dmesg /etc/rc.d/gptboot /etc/rc.d/hostapd /etc/rc.d/mdconfig2 /etc/rc.d/motd /etc/rc.d/newsyslog /etc/rc.d/os-release /etc/rc.d/virecover /etc/rc.d/wpa_supplicant 
    30	/etc/rc.d/syslogd 
    31	/etc/rc.d/auditd /etc/rc.d/bsnmpd /etc/rc.d/hastd /etc/rc.d/ntpdate /etc/rc.d/power_profile /etc/rc.d/pwcheck /etc/rc.d/savecore /etc/rc.d/watchdogd 
```

In result, rc-messages prior to `syslogd(8)` are only availalbe in console. The
first ones are from `auditd` service in `/var/log/messages`:

```
Aug 29 14:40:08 nuc kernel: em0: link state changed to UP
Aug 29 14:40:08 nuc root[26487]: /etc/rc: DEBUG: checkyesno: auditd_enable is set to NO.
```

Configure `syslog(3)` to log all console messages:

```console
% grep console.info /etc/syslog.conf 
console.info					/var/log/console.log
```

## SSH

```console
% sysrc -s sshd -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/sshd /usr/local/etc/rc.conf.d/sshd
# sysrc -x sshd_enable
# sysrc -f /etc/rc.conf.d/sshd sshd_enable=yes
# service sshd restart
```

By default, SSH server listens on all IP addresses. Restrict it to the host IP
on LAN `em0` and WLAN `wlan0` and only allow operator user `op` to connect:

```console
# sockstat -4 | grep sshd
root     sshd          21 7   tcp4   *:22                  *:*
# sysrc -f /etc/rc.conf.d/sshd sshd_flags="-o ListenAddress=192.168.1.100 -o ListenAddress=192.168.1.101 -o AllowUsers=op"
# service sshd restart
# sockstat -4 | grep sshd
root     sshd       32886 7   tcp4   192.168.1.101:22      *:*
root     sshd       32886 8   tcp4   192.168.1.100:22      *:*
```

## ZFS

```console
% sysrc -s zfs -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/zfs /usr/local/etc/rc.conf.d/zfs
# sysrc -x zfs_enable
# sysrc -f /etc/rc.conf.d/zfs zfs_enable=yes
# shutdown -r now
```

# APPLICATIONS

Ref: https://docs.freebsd.org/en/books/handbook/ports/

Ports build software from the source. Packages are pre-built binaries. There
might be multiple packages for the same port representing the same application
with different configuraiton options. Not every port has a binary package.

Switch from Quarterly to Latest packages:

```console
# mkdir -p /usr/local/etc/pkg/repos
# cat <<eof > /usr/local/etc/pkg/repos/FreeBSD.conf 
FreeBSD: { url: "pkg+http://pkg.FreeBSD.org/${ABI}/latest" }
eof
# pkg update -f
```

List of installed packages:

```console
% pkg prime-list
doas
gpu-firmware-intel-kmod-skylake
pkg
rsync
vim-tiny
wifi-firmware-iwlwifi-kmod-8000
```

## Doas

Let the operator monitor the system:

```console
# pkg install doas
# cat <<eof >/usr/local/etc/doas.conf
permit nopass op cmd sockstat
eof
```

## Vim

There are different pre-built configurations of
[vim](https://www.freshports.org/editors/vim/). `vim-tiny` only includes binary
without runtime files. It makes `vim-tiny` suitable for minimal installations.

```console
# pkg install vim-tiny
```
