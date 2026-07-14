<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**ssl** - manage SSL certificates

# DESCRIPTION

The instructions below create a certificate autohority (CA) and use it to sign
a certificate for LDAP service. It uses different hosts to create and use
certificates.

Keep in mind that trusted certificates:

- must reside under `/usr/local/share/certs/` (see `TRUSTPATH` in certctl(8)).
- installed certificates live under `/etc/ssl/certs/`

Use `certctl rehash` to install new certificates.

```
# mkdir -vp /usr/local/share/certs
# mkdir -vp /usr/local/share/ssl/{certs,private}
# chmod -v 700 /usr/local/share/ssl/private
```

## Create certificates

Generate a CA certificate:

```
# cd /usr/local/share/ssl/private

-- the certificate is valid for 5y (1825 days)
# openssl req \
    -days 1825 \
    -keyout ca.key \
    -new \
    -nodes \
    -out /usr/local/share/certs/ca.crt \
    -x509
... # fqdn: ca.ssl.lab.net

# certctl rehash
```

Create an LDAP certificate:

```
# cd /usr/local/share/ssl/private

-- the certificate is valid for 1y (365 days)
-- .csr stands for "certificate signing request"
# openssl req \
    -days 365 \
    -keyout ldap.key \
    -new \
    -nodes \
    -out ldap.csr
... # fqdn: ldap.lab.net

# openssl x509 \
    -CA /usr/local/share/certs/ca.crt \
    -CAkey ca.key \
    -days 365 \
    -in ldap.csr \
    -out /usr/local/share/ssl/certs/ldap.crt \
    -req
```

## Install certificates

To use the certificate in the LDAP server:

```
# mkdir -vp /usr/local/share/ssl/{certs,private}
# chmod -v 750 /usr/local/share/ssl/private
```

Install CA certificate:

```
# mkdir /usr/local/share/certs

-- copy ca.crt to /usr/local/share/ssl/certs

# certctl rehash
```

Install ldap certificate:

- copy ldap.cert to /usr/local/share/ssl/certs
- copy ldap.key to /usr/local/share/ssl/private
- make sure the process user, e.g. `:ldap`, can read the key and /usr/local/share/ssl/private

Validate

```
# openssl verify -verbose /usr/local/share/ssl/certs/ldap.crt
/usr/local/share/ssl/certs/ldap.crt: OK
```
