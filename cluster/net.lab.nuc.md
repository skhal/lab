# NAME

**nuc.lab.net** - FreeBSD cluster on Intel NUC Kit NUC6i7KYK.


# DESCRIPTION

`nuc.lab.net` is a host running FreeBSD cluster on Intel NUC Kit NUC6i7KYK:

  * 32 GB of RAM
  * 2x 1TB M.2 NVMe SSD

```console
% pciconf -lv | grep -B 4 'subclass.*NVM'
nvme0@pci0:61:0:0:	class=0x010802 rev=0x00 hdr=0x00 vendor=0x15b7 device=0x5006 subvendor=0x15b7 subdevice=0x5006
    vendor     = 'Sandisk Corp'
    device     = 'SanDisk Extreme Pro / WD Black SN750 / PC SN730 / Red SN700 NVMe SSD'
    class      = mass storage
    subclass   = NVM
nvme1@pci0:62:0:0:	class=0x010802 rev=0x00 hdr=0x00 vendor=0x15b7 device=0x5006 subvendor=0x15b7 subdevice=0x5006
    vendor     = 'Sandisk Corp'
    device     = 'SanDisk Extreme Pro / WD Black SN750 / PC SN730 / Red SN700 NVMe SSD'
    class      = mass storage
    subclass   = NVM
```

## Install FreeBSD

Ref: https://docs.freebsd.org/en/books/handbook/bsdinstall/

Download AMD64 image and verify the checksum. MacOS does not include `xz(1)`
in the base installation, therefore use `.img`:

```console
% curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/FreeBSD-14.3-RELEASE-amd64-mini-memstick.img
% curl -O https://download.freebsd.org/releases/amd64/amd64/ISO-IMAGES/14.3/CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64
% shasum \
    -c CHECKSUM.SHA512-FreeBSD-14.3-RELEASE-amd64 \
    --ignore-missing
FreeBSD-14.3-RELEASE-amd64-mini-memstick.img: OK
```

Write the image to the memstick:

```console
% disktuil list
...
/dev/disk6 (external, physical):
#:                       TYPE NAME                    SIZE       IDENTIFIER
0:     FDisk_partition_scheme                        *7.9 GB     disk6
1:                       0xEF                         34.1 MB    disk6s1
2:                    FreeBSD                         1.4 GB     disk6s2
	      (free space)                         6.5 GB     -
# diskutil umountDisk /dev/disk6
Unmount of all volumes on disk6 was successful
# dd \
    status=progress \
    if=FreeBSD-14.3-RELEASE-amd64-mini-memstick.img \
    of=/dev/disk6 \
    bs=1m \
    conv=sync
532676608 bytes (533 MB, 508 MiB) transferred 22.021s, 24 MB/s
508+0 records in
508+0 records out
532676608 bytes transferred in 22.059363 secs (24147416 bytes/sec)
% diskutil eject /dev/disk6
Disk /dev/disk6 ejected
```

Boot from the memstick and follow the instructions:

  * Keyboard: default layout (US)
  * Hostname: `nuc.lab.net`
  * Network: `em0` IPv4 DHCP
  * Filesystem: ZFS with mirrored pool
  * Enable services:
    - `dumpdev` to dump kernel crashes to `/var/crash`
    - `ntpd` for clock synchronization
    - `ntpd_sync_on_start` to sync time on `ntpd` start
    - `sshd` for SSH server
  * Turn on security hardening:
    - `hide_uids` hide processes as other users
    - `hide_gids` hide processes as other groups
    - `hide_jail` hide processes in jails
    - `read_msgbuf` no kernel msgbuf read for unprivileged
    - `proc_debug` no proc debug for unprivileged
    - `random_pid` random PID for new processes
    - `clear_tmp` clean `/tmp` on startup
    - `disable_syslogd` no syslogd network socket
    - `secure_console` console password prompt
  * Install FreeBSD handbook

## Update the OS

Ref: https://docs.freebsd.org/en/books/handbook/cutting-edge/

Check the running version of installed kernel `-k`, running kernel `-r`,
and running userland `-u`:

```console
% freebsd-version -kru
14.3-RELEASE
14.3-RELEASE
14.3-RELEASE
```

Fetch and install updates. The system will auto-reboot if there is
a kernel update, otherwise it restarts the updated services only.
It is still a good idea to restart the node.

```console
# freebsd-update fetch
# freebsd-update install
# reboot
```

Verify the running version is up to date, all three should match:

```console
% freebsd-version -kru
14.3-RELEASE-p2
14.3-RELEASE-p2
14.3-RELEASE-p2
```

## Create users

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

Create a group to manage jails:

```console
# pw groupadd -g 1001 -n jail
```

## Configure services

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

### Hostname

```console
% sysrc -s hostname -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/hostname /usr/local/etc/rc.conf.d/hostname
# sysrc -x hostname
# sysrc -f /etc/rc.conf.d/hostname hostname="nuc.lab.net"
# service hostname restart
```

### Jail

```console
# sysrc -f /etc/rc.conf.d/jail jail_enable="YES"
jail_enable: NO -> YES
# sysrc -f /etc/rc.conf.d/jail jail_parallel_start="YES"
jail_parallel_start: NO -> YES
```

### Moused

```console
% sysrc -s moused -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/moused /usr/local/etc/rc.conf.d/moused
# sysrc -x moused_nondefault_enable
# sysrc -f /etc/rc.conf.d/moused moused_nondefault_enable=no
moused_nondefault_enable: YES -> no
```

### Network

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

Wireless: https://docs.freebsd.org/en/books/handbook/network/#network-wireless

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

### Network Time server

```console
% sysrc -s ntpd -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/ntpd /usr/local/etc/rc.conf.d/ntpd
# sysrc -x ntpd_sync_on_start
# sysrc -x ntpd_enable
# sysrc -f /etc/rc.conf.d/ntpd ntpd_enable=yes
# sysrc -f /etc/rc.conf.d/ntpd ntpd_sync_on_start=yes
# service ntpd restart
```

### Routing

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

### RC

Enable rc debug and info logging:

```console
# sysrc rc_debug=yes
# sysrc rc_info=yes
```

### Syslog

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

### SSH

```console
% sysrc -s sshd -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/sshd /usr/local/etc/rc.conf.d/sshd
# sysrc -x sshd_enable
# sysrc -f /etc/rc.conf.d/sshd sshd_enable=yes
# service sshd restart
```

By default, SSH server listens on all IP addresses. Restrict it to the host IP:

```console
# sockstat -4 | grep sshd
root     sshd          21 7   tcp4   *:22                  *:*
# sysrc -f /etc/rc.conf.d/sshd sshd_flags+="-o ListenAddress=192.168.1.100"
# service sshd restart
# sockstat -4 | grep sshd
root     sshd          21 7   tcp4   192.168.1.100:22      *:*
```

Allow only user `op` to connect over SSH:

```console
# grep sshd_flags /etc/rc.conf.d/sshd
sshd_flags="-o ListenAddress=192.168.1.100 -o AllowUsers=op"
```

Optional: make SSH server to listen to connection requests on wireless network:

```console
# sysrc -f /etc/rc.conf.d/sshd sshd_flags+=" -o ListenAddress=192.168.1.101"
sshd_flags: -o ListenAddress=192.168.1.100 -> -o ListenAddress=192.168.1.100 -o ListenAddress=192.168.1.101
# service sshd restart
```

### ZFS

```console
% sysrc -s zfs -l
/etc/rc.conf /etc/rc.conf.local /etc/rc.conf.d/zfs /usr/local/etc/rc.conf.d/zfs
# sysrc -x zfs_enable
# sysrc -f /etc/rc.conf.d/zfs zfs_enable=yes
# shutdown -r now
```

## Applications

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

Install packages:

```console
% pkg prime-list
doas
gpu-firmware-intel-kmod-skylake
pkg
vim-tiny
wifi-firmware-iwlwifi-kmod-8000
```

### Doas

Let members of `jail` group manage jails:

```console
# cat <<eof >/usr/local/etc/doas.conf
permit nopass :jail cmd jail
permit nopass :jail cmd jexec
eof
```
