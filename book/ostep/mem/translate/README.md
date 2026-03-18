<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**translate** - map between virtual and physical memory address

# SYNOPSIS

```
translate
```

# DESCRIPTION

`translate` demonstrates address translation between virtual and physical
address space using interpolation.

Basic interpolation uses two registers:

- *base*: the offset of the memory address.

- *bound*: the size of the address block to restrict access outside process'
  memory boundaries.

The physical address is:

```
phys address = virt address + base
```
