<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**container** - largest container

PROBLEM
=======

Consider an array of heights {hi}. A container is a box, formed from two heights at indices i and j and respective heights hi and hj. It's volume is defined as:

```
min(hi, hj) * (j - i)
```

Find a box with maximum volume and return the volume value.

EXAMPLE
=======

*Input*: [1, 2, 3]

*Output*: 2

SEE ALSO
========

[Solution](./solution.md)
