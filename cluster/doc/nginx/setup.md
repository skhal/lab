Name
====

**nginx** - Nginx server setup

Description
===========

```console
# pkg install nginx
```

Pull configuration to serve static user content from ~/www at nginx.example.com/~user/:

```console
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/cluster/doc/nginx/data/nginx.conf.diff
# patch -b -i /tmp/nginx.conf.diff /usr/local/etc/nginx/nginx.conf
```

Enable the server:

```console
# mkdir /usr/local/etc/rc.conf.d
# sysrc -f /usr/local/etc/rc.conf.d/nginx nginx_enable=yes
```
