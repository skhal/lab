# NAME

**rcorder** -- list `rc(8)` scripts in dependency ordering.


# DESCRIPTION

`rcorder(8)` prints a list of scripts in dependency ordering. Every script
must have dependencies set in a comment block of the form:

```sh
# BEFORE: <before>
# KEYWORD: <keyword>
# PROVIDE: <provide>
# REQUIRE: <require>
```

The values `<before>`, `<keyword>`, `<provide>`, and `<require>` are space
separated tokens.

The `KEYWORD` is a tag (see `rc(8)`). `rcorder` can be forced to skip (`-s tag`)
or keep (`-k tag`) scripts with a tag.

Examples:

  * `-s nojail`: [skip](https://github.com/freebsd/freebsd-src/blob/9bfb1405332c6c847dd29e4db4dd3afb56662021/libexec/rc/rc#L87)
    services that are not meant to run inside a jail. ZFS daemon can't run in
    a jail  
    ```console
    % grep '\bnojail\b' ./libexec/rc/rc.d/* | grep zfs
    ./libexec/rc/rc.d/zfsd:# KEYWORD: nojail shutdown
    ```

  * `-s nojailvnet`: [skip](https://github.com/freebsd/freebsd-src/blob/9bfb1405332c6c847dd29e4db4dd3afb56662021/libexec/rc/rc#L89)
    services in the jails without vnet. `netif` can't run in a jail without a
    virtual network  
    ```console
    % grep '\bnojailvnet\b' ./libexec/rc/rc.d/* | grep netif
    ./libexec/rc/rc.d/netif:# KEYWORD: nojailvnet
    ```

One way to detect a jail is to use `sysctl(8)` to check the flag (0 - no,
1 - yes):

```console
% sysctl -d security.jail.jailed
security.jail.jailed: Process in jail?
% sysctl -d security.jail.vnet
security.jail.vnet: Jail owns vnet?
```

The following example generates a dependency ordering of the services included
in the base installation (see [`rcorder_base.txt`](./rcorder_base.txt)):

```console
% rcorder /etc/rc.d/* | head -3
/etc/rc.d/dhclient
/etc/rc.d/dumpon
/etc/rc.d/dnctl
```

The list does not include service that are installed with ports or packages in
`/usr/local/etc/rc.d/`. Include these services in the `rcorder(8)` arguments:

```console
% rcorder /etc/rc.d/* | wc -l
    173
% rcorder /etc/rc.d/* /usr/local/etc/rc.d/* | wc -l
    175
```

One caveat is that `rcorder(8)` does not check whether a service is enabled.
Use `service(8)` to list only enabled services:

```console
% service -r | wc -l
    96
```
