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
