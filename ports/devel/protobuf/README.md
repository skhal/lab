<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**devel/protobuf** -- protocol buffers for FreeBSD

# DESCRIPTION

As of Sep '26, the official FreeBSD port
[`devel/protobuf`](https://github.com/freebsd/freebsd-ports/tree/8a2e817ddae3c8ecbb4d30869f753f6ea53ee43d/devel/protobuf)
is based on Protobuf version
[v29.6](https://github.com/protocolbuffers/protobuf/releases/tag/v29.6)
and does not support 2024 editions, added in
[v32.0](https://github.com/protocolbuffers/protobuf/releases/tag/v32.0).

This port is a patched clone of the official one. It bumps Protobuf version
to at least v32.0.

## Why not to patch the official port?

We communicated with FreeBSD port maintainer in May '26. It turns out that some
of the dependent ports are not compatible with this change, e.g.
`devel/protobuf-c`.

There is a known issue that prevents protobuf-c from moving forward. See pull
request
[PR#797](https://github.com/protobuf-c/protobuf-c/pull/797),
created in Mar '26.
