<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**mem** - virtual memory

# DESCRIPTION

- [Layout](./layout/): print address of the code block, stack, and heap.
- [Address Translation](./translate/): map between virtual and physical
  address space.
- [Segmented Address Translation](./segment/): address translation with support
  for memory segments in the virtual address space.
- [Memory Allocator](./allocator/): malloc(3)-like memory allocator.
