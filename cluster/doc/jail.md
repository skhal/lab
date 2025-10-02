# NAME

**jail** - setup jails on `nuc.lab.net`


# NO NET JAIL

The new jail `dev` is isolated and does not have a network access:

```console
# zfs clone zroot/jail/template/14.3-RELEASE@p1 zroot/jail/container/dev
# cat /etc/jail.conf.d/dev.conf
dev {
  # empty
}
```

# VNET JAIL

Ref: https://freebsdfoundation.org/wp-content/uploads/2020/03/Jail-vnet-by-Examples.pdf

Virtual Network (VNET) adds a networking stack to the jail, isolated from the
host system. It includes interfaces, addresses, routing tables and firewall
rules.

The setup uses two network interfaces to connect jails using bridge(4):

* `bridge0` connects jails to the Internet. Subnet: `192.168.1.0/24`
* `bridge1` isolates intra-jail traffic. Subnet: `10.0.1.0/24`

Create bridges and give `bridge0` Internet access:

```console
# sysrc -f /etc/rc.conf.d/network cloned_interfaces+="bridge0 bridge1"
cloned_interfaces:  -> bridge0 bridge1
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge0="addm em0 up descr jail:em"
ifconfig_bridge0:  -> addm em0 up descr jail:em
# sysrc -f /etc/rc.conf.d/network ifconfig_bridge1="up descr jail:lo"
ifconfig_bridge0:  -> addm em0 up descr jail:lo
```

> [!Note]
> The bridge and both ends of epair must be in the UP state for the packets to
> flow.

Let [`rc.jail`](https://github.com/skhal/lab/blob/main/freebsd/rc/rc.jail)
script managet jail's epair(4), connect it to the bridge, and assign the IP
address.

```console
% cat /etc/jail.conf.d/dev.conf
dev {
  $id = "110";

  $bridge0 = "bridge0";
  $bridge0_ip = "${bridge0}:192.168.1.${id}/24";

  $bridge1 = "bridge1";
  $bridge1_ip = "${bridge1}:10.0.1.${id}/24";

  $bridges = "${bridge0} ${bridge1}";
  $bridgeips = "${bridge0_ip} ${bridge1_ip}";

  vnet;
  devfs_ruleset = 5;

  exec.prestart = "/bin/sh /usr/local/etc/rc.jail prestart ${name} ${bridges}";
  exec.created  = "/bin/sh /usr/local/etc/rc.jail created ${name} ${bridgeips}";
  exec.poststop = "/bin/sh /usr/local/etc/rc.jail poststop ${name} ${bridges}";
}
```
