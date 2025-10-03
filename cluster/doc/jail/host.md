# NAME

**host** - host setup for jails


# DESCRIPTION

## Directory tree

Jails reside under `/jail`. There are three folders:
  * `/jail/image` stores downloaded user lands in compressed format.
  * `/jail/template` holds base templates to create jails.
  * `/jail/container` keeps running jail.

```console
# zfs create -o mountpoint=/jail zroot/jail
# zfs create zroot/jail/image
# zfs create zroot/jail/template
# zfs create zroot/jail/container
```

## Permissions

Create `jail` user group:

```console
# pw groupadd -g 1001 -n jail
```

`jail` members can manage and enter jails:

```console
# cat <<eof >> /usr/local/etc/doas.conf
permit nopass :jail cmd jail
permit nopass :jail cmd jexec
eof
```

Members of `jail` can manage jail datasets:

```console
# zfs allow -s @mount  mount,canmount,mountpoint zroot/jail
# zfs allow -s @create create,destroy,@mount zroot/jail
# zfs allow -g jail @mount,@create,readonly zroot/jail
```

Add system operator to the `jail` group:

```console
# pw groupmod -m op -n jail
```

## Jail service

`jail(8)` reads configuration from `/etc/jail.conf`. Set default configuration
parameters and pick up individual jail configurations from `/etc/jail.conf.d/`

```console
# mkdir /etc/jail.conf.d
# cat /etc/jail.conf
host.hostname = "${name}.lab.net";

path = "/jail/container/${name}";
exec.consolelog = "/var/log/jail_${name}.log";

exec.start = "/bin/sh /etc/rc";
exec.stop = "/bin/sh /etc/rc.shutdown";

.include "/etc/jail.conf.d/*.conf";
```

Enable `jail` service and stop jails in the reverse order to ensure dependencies
are satisfied:

```console
# sysrc -f /etc/rc.conf.d/jail jail_reverse_stop=yes
jail_reverse_stop: NO -> yes
# sysrc -f /etc/rc.conf.d/jail jail_enable=yes
jail_enable: NO -> yes
```


# SEE ALSO

  * https://docs.freebsd.org/en/books/handbook/jails/
