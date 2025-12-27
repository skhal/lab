Name
====

**jail** -- notes on FreeBSD jails

Notes
=====

Jail parameters
---------------

We assume one runs a VNET dummy jail - [demo.conf](./data/etc/jail.conf.d/demo.conf).

Full list of supported parameters (the command is equivalent to `jls -nj demo all` - [usr.sbin/jls/jls.c](https://github.com/freebsd/freebsd-src/blob/086bedb11a853801e82234b8a1a64f0df52d9e52/usr.sbin/jls/jls.c#L179)\):

```console
% % jls -nj demo | tr ' ' '\n' | sort | head -n 3
allow.mount.nodevfs
allow.mount.nofdescfs
allow.mount.nolinprocfs
```

See [jflags_demo](./data/jflags_demo).

JSON output
-----------

jls(8) integrates with libxo(3) that can output data in JSON format:

```console
% jls --libxo=json -j demo jid name | jq
{
  "__version": "2",
  "jail-information": {
    "jail": [
      {
        "jid": "8",
        "name": "demo"
      }
    ]
  }
}
```
