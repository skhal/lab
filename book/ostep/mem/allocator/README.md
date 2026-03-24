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

- *backward*: coalesce consecutive previous blocks staring from the released
  block.

- *bidi*: bidirectional coalesce of consecutive free blocks, moving forward
  and backward.

# EXAMPLE

```console
% allocator -base 1024 -size 2048 -n 5 -c bidi -n 15
configuration:
  base: 1024 size: 2048 coalesce: bidi
  [1] blocks 2046:1026[-F]

trace:
  malloc(1192)
    [1] addresses 1026
    [2] blocks 1192:1026[-A] 852:2220[PF]
  malloc(171)
    [2] addresses 1026 2220
    [3] blocks 1192:1026[-A] 171:2220[PA] 679:2393[PF]
  malloc(800) malloc(800): insufficient memory
  malloc(1785) malloc(1785): insufficient memory
  malloc(591)
    [3] addresses 1026 2220 2393
    [4] blocks 1192:1026[-A] 171:2220[PA] 591:2393[PA] 86:2986[PF]
  malloc(1275) malloc(1275): insufficient memory
  malloc(586) malloc(586): insufficient memory
  free(1026)
    [2] addresses 2220 2393
    [4] blocks 1192:1026[-F] 171:2220[-A] 591:2393[PA] 86:2986[PF]
  malloc(1533) malloc(1533): insufficient memory
  malloc(265)
    [3] addresses 2220 2393 1026
    [5] blocks 265:1026[-A] 925:1293[PF] 171:2220[-A] 591:2393[PA] 86:2986[PF]
  malloc(920)
    [4] addresses 2220 2393 1026 1293
    [6] blocks 265:1026[-A] 920:1293[PA] 3:2215[PF] 171:2220[-A] 591:2393[PA] 86:2986[PF]
  free(1026)
    [3] addresses 2220 2393 1293
    [6] blocks 265:1026[-F] 920:1293[-A] 3:2215[PF] 171:2220[-A] 591:2393[PA] 86:2986[PF]
  free(2220)
    [2] addresses 2393 1293
    [5] blocks 265:1026[-F] 920:1293[-A] 176:2215[PF] 591:2393[-A] 86:2986[PF]
  free(2393)
    [1] addresses 1293
    [3] blocks 265:1026[-F] 920:1293[-A] 857:2215[PF]
  free(1293)
    [0] addresses
    [1] blocks 2046:1026[-F]
```
