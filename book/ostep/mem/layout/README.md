<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**layout** - show address of the code block, stack, and heap.

# SYNOPSIS

```
layout
```

# DESCRIPTION

`layout` demonstrates memory layout of core blocks: code, stack, heap.

Traditionally, stack goes first. The stack and heap grow from the opposite
ends to guarantee dynamic nature of the two (it is impossible to predict how
large each is going to be).

Possible output:

```console
% layout
code:   main() at 0x201670
stack:  argc at 0x82080efb8
heap:   malloc() at 0x361343812000
```
