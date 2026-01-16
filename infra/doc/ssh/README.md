<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**ssh** - secure shell

Jump Host
=========

A jump host allows clients to connect to the servers on a different network, sitting behind a proxy server.

Consider following topology:

```
[client 192.168.1.10/24]
  |
  | Network-A: 192.168.1.0/24
  |
[proxy 192.168.1.11/24]
  |
  | Network-B: 10.0.0.0/24
  |
[server 10.0.0.1/24]
```

The client does not know about the Network-B or any of the servers connected to it but would like to make an SSH connection to the server. A proxy jump is at rescue:

```console
% ssh -J foo@proxy bar@server
```

The above ssh command first connects to the proxy, and then establish a forwarding TCP connection to the server from the proxy. Keep in mind that DNS resolution of the server happens on the proxy node, i.e. a DNS record of the server to 10.0.0.1 works fine.

Use authorization keys to automate the process with an SSH configuration on the client:

```
Host proxy
  HostHame 192.168.1.11
  IdentityFile ~/.ssh/id_proxy

Host server
  ProxyJump foo@proxy
```
