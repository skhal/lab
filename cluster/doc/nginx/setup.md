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
# patch -i /tmp/nginx.conf.diff /usr/local/etc/nginx/nginx.conf
```
