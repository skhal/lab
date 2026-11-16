<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Pkg

Pkg configures pkg(8) with latest branch, adds pkg aliases, and installed
requested packages.

# Example

Add pkg(8) alias `download` and install python3 package:

```yaml
- role: pkg
  vars:
    pkg_aliases:
      - name: download
        value: fetch
    pkg_install:
      - python3
```
