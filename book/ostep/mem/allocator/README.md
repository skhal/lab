<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**allocator** - malloc(3) like memory allocator

# DESCRIPTION

Allocator simulates malloc(3) like memory management with minimal API:

```go
type Heap interface {
    Malloc(size int) (address int, err error)
    Free(address int)
}
```

It manages a pre-allocated block of bytes to simulate memory obtained from the
Kernel via mmap(2) or brk(2).

## Free List

The allocator uses implicit free list. The managed memory consists of blocks
with *header* and *payload*. The header holds block size and flags to indicate
whether the block is allocated or free. The payload is actual address space
returned to the use from Malloc().

## Coalesce

The allocator supports following strategies to coalesce free spacea after
Free() call to reduce fragmentation:

- *noop*: disable coalesce.

- *forward*: only coalesce next blocks started from the released block.

# EXAMPLE

```console
% allocator -base 1024 -size 2048 -n 5 -c forward
configuration:
  base: 1024 size: 2048 coalesce: forward
  [1] free blocks 2046:1026

trace:
  malloc(225)
    [1] allocations 1026
    [1] free blocks 1819:1253
  malloc(326)
    [2] allocations 1026 1253
    [1] free blocks 1491:1581
  free(1026)
    [1] allocations 1253
    [2] free blocks 225:1026 1491:1581
  malloc(1892) malloc(1892): insufficient memory
  free(1253)
    [0] allocations
    [2] free blocks 225:1026 1819:1253
```
