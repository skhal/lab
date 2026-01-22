<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**feedreader** - CLI streaming feeds reader.

Synopsis
========

```console
feedreader -f file
```

Description
===========

Feedreader streams feeds from sources, configured in the `-f file` flag. It could be a local feed file or an online feed (under development).

The feeds configuration is in [text](https://protobuf.dev/reference/protobuf/textformat-spec/) Protobuf format (see [schema](../../internal/pb/feed.proto)\):

```txtpb
feeds {
  name: "demo"
  source {
    file: "~/rss_feed.xml"
    kind: KIND_RSS
  }
}
```
