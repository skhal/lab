Name
====

**jail** -- notes on FreeBSD jails

Notes
=====

We assume one runs a VNET dummy jail - [demo.conf](./data/etc/jail.conf.d/demo.conf).

Full list of supported parameters (the command is equivalent to `jls -nj demO all`\):

```console
% % jls -nj demo | tr ' ' '\n' | sort | head -n 3
allow.mount.nodevfs
allow.mount.nofdescfs
allow.mount.nolinprocfs
```

See [jflags_demo](./data/jflags_demo).
