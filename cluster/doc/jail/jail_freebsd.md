# NAME

**jail-freebsd** - create a FreeBSD jail


# DESCRIPTION

## ZFS Bootstrap

Start from FreeBSD template:

```console
% zfs list -t snapshot zroot/jail/template/14.3-RELEASE
NAME                                    USED  AVAIL  REFER  MOUNTPOINT
zroot/jail/template/14.3-RELEASE@p2     144K      -   459M  -
zroot/jail/template/14.3-RELEASE@p2.1   152K      -   459M  -
zroot/jail/template/14.3-RELEASE@p2.2   128K      -   459M  -
zroot/jail/template/14.3-RELEASE@p2.3   120K      -   459M  -
zroot/jail/template/14.3-RELEASE@p2.4     8K      -   459M  -
# zfs clone zroot/jail/template/14.3-RELEASE@p2.4 zroot/jail/container/demo
```

## Configuration

Use `rc.jail` script on the host environment manage jail epair(4) for every
bridge:

```console
# cat /etc/jail.conf.d/demo.conf
demo {
  $id = "90";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  host.hostname = "${name}.lab.net";
  path = "/jail/container/${name}";

  vnet;
  allow.raw_sockets;

  mount.devfs;
  devfs_ruleset = 5;

  enforce_statfs = 1;

  exec.clean;
  exec.consolelog = "/var/log/jail_${name}.log";

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.start    = "/bin/sh /etc/rc";

  exec.stop     = "/bin/sh /etc/rc.shutdown";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```

## Start jail

Manually:

```console
# service jail start demo
```

Auto-start:

```console
# sysrc -f /etc/rc.conf.d/jail jail_list+=demo
jail_list: '' -> demo
```


# SEE ALSO

* https://freebsdfoundation.org/wp-content/uploads/2020/03/Jail-vnet-by-Examples.pdf
