<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

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
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/nginx/data/nginx.conf.diff
# patch -lb -i /tmp/nginx.conf.diff /usr/local/etc/nginx/nginx.conf
```

Nginx configuration uses syslogd(8) for logging. Configure it:

```console
# mkdir -v /usr/local/etc/syslog.d
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/nginx/data/syslog_nginx.conf.diff
# mv /tmp/syslog_nginx.conf /usr/local/etc/syslog.d/nginx.conf
```

Rotate the logs with newsyslog(8):

```console
# mkdir -v /usr/local/etc/newsyslog.conf.d
# fetch -o /tmp https://raw.githubusercontent.com/skhal/lab/refs/heads/main/infra/doc/nginx/data/newsyslog_nginx.conf.diff
# mv /tmp/newsyslog_nginx.conf /usr/local/etc/newsyslog.conf.d/nginx.conf
```

Enable the server:

```console
# mkdir /usr/local/etc/rc.conf.d
# sysrc -f /usr/local/etc/rc.conf.d/nginx nginx_enable=yes
```
