<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**vpn** - Virtual Private Network

# DEDSCRIPTION

Virtual Private Network (VPN) allows to connect to private networks from
anywhere.

**WARNING**: keep private network on a different sub-network compared to local
network that provides connection. For example, local connection may use network
192.168.1.0/24. Use different remote private network, say 192.168.2.0/24. This
way VPN can separate traffic between the networks.

## Dynamic DNS

Most of Internet Service Providers (ISP) assign dynamic IPs to clients, meaning
that external interface IP rotates once in a wile. Dynamic DNS allows to
overcome the issue by updating a DNS entry every time external IP changes.

UniFI cloud gateways support Dynamic DNS. It is available under: Settings /
Network / ISP / Advanced (Manual) / Dynamic DNS.

Setup DDNS with noip ([help.ui.com](https://help.ui.com/hc/en-us/articles/9203184738583-UniFi-Gateway-Dynamic-DNS)).

## VPN Server

UniFI cloud gateway UI: Settings / VPN / VPN Server.

Create a WireGuard server ([help.ui.com](https://help.ui.com/hc/en-us/articles/115005445768-UniFi-Gateway-WireGuard-VPN-Server)). The instructions are straightforward,
no additional setup is necessary.

## VPN Client

Create a client in UniFI cloud gateway UI under VPN Server. Download client
settings.

Install WireGuard client on the client machine and import the settings to create
a tunnel.

At this point the tunnel will route all traffic through VPN. Test it to make
sure VPN operates as expected.

Update tunnel configuration to use DDNS instead of IP to connect to VPN server:

```
Endpoint = sklab.ddns.net:51820
```

## VPN Split Tunnel

Instead of routing all VPN traffic through VPN server, it is beneficial to only
route traffic that is dedicated for the private network.

Change settings on the WireGuard client, in the tunnel management settings:

```
AllowedIPs = 192.168.10.0/24, 192.168.2.0/24
```

The configuration above routes only traffic for 192.168.10 and 192.168.2
sub-networks through VPN. Everything else uses local connection, e.g. surfing
web, video, music, etc.
