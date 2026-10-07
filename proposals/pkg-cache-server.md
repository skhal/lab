<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

## Background

Nginx is a caching PKG server. It is setup with a single virtual server to serve different port collections:

- FreeBSD base - base system
- FreeBSD kmods - Kernel modules
- FreeBSD ports - user ports

Clients access the server via `https://pkg.lab.net`, with a single TLS certificate.

https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/templates/usr/local/etc/nginx/nginx.conf.j2#L27-L38

The server uses a URL path to distinguish port collections. For example, FreeBSD ports collection uses `/freebsd/ports/...` path prefix:

https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/files/usr/local/etc/nginx/conf.d/http/server/freebsd-ports.conf#L6-L8

The server rewrites the URL to proxy the request to `http://pkg.freebsd.org`:

https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/files/usr/local/etc/nginx/conf.d/http/server/freebsd-ports.conf#L16

## Problem

Nginx configuration uses include-statements to load different configurations for each port collection.

Each collection config has two components:

1. `proxy_cache_path` to configure cache, i.e. location on the filesystem, expiration policy, etc.
   https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/files/usr/local/etc/nginx/conf.d/http/freebsd-ports-proxy-cache.conf#L6-L7

2. `location /.../ { ... }` to serve the ports collection based on the URL path and rewrite the URL to proxy the request to FreeBSD official repository:
   https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/files/usr/local/etc/nginx/conf.d/http/server/freebsd-ports.conf#L6-L19

This complicates the main `nginx.conf` in the following ways:

- the configuration must explicitly include `http.proxy_cache_path` constructs:
  https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/templates/usr/local/etc/nginx/nginx.conf.j2#L23-L27

- there is another include for `http.server.location`:
  https://github.com/skhal/lab/blob/68f3ed358e1347af88b1aecf4b949ef3858a8a4d/infra/ansible/roles/local_pkg/templates/usr/local/etc/nginx/nginx.conf.j2#L40-L41

In result, the included configuration must be carefully named or placed on the file system to distinguish the two. Current setup takes the following approach:

```
/usr/local/etc/nginx/conf.d/
  http/                             « http-block statements
    freebsd-ports-proxy-cache.conf  « proxy_cache_path for FreeBSD ports collection
    server/                         « http.server-block statements like location
      freebsd-ports.conf            « location block for FreeBSD ports collection
```

## Proposal

An alternative solution to the problem is to separate port collections by
virtual servers instead of location. This change simplifies main Nginx
configuration file:

```
# file: /usr/local/etc/nginx/nginx.conf
http {
  include /usr/local/etc/nginx/conf.d/*;
}
```

With each port collection's configuration encapsulate all the parameters,
including http- and http.server-block configs:

```
# file: /usr/local/etc/nginx/conf.d/freebsd-ports.conf
proxy_cache_path /var/cache/pkg.freebsd.ports levels=1:2 use_temp_path=off
                 keys_zone=freebsd-ports:10m inactive=14d;
server {
  listen :443 ssl;
  server_name freebsd-ports.pkg.lab.net;

  ssl_certificate /usr/locl/etc/nginx/ssl/freebsd-ports.crt;
  ssl_certificate_key /usr/locl/etc/nginx/ssl/private/freebsd-ports.key;

  location / {
    proxy_cache freebsd-ports;
    proxy_pass  https://pkg.FreeBSD.org;
    add_header	X-Cache-Status $upstream_cache_status;
  }
}
```

Keep in mind that the new location configuration does not perform URL rewrites
in this setup.

This change requires coordination with additional changes:

- **DNS**: add an Alias (CNAME) and a Service (SRV) records to DNS server:

  ```
  freebsd-ports.pkg IN CNAME pkg
  _https._tcp.freebsd-ports.pkg.lab.net. IN SRV 10 100 443 freebsd-ports.pkg.lab.net.
  ```

- **Certificate**: add Subject Alternative Name (SAN) to pkg.lab.net TLS
  certificate to serve sub-domains, e.g. `freebsd-ports.pkg.lab.net`.
