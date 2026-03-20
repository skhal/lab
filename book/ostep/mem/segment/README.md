<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**segment** - segmented address translation

# DESCRIPTION

At the very minimum, virtual address space consists of code, stack and heap,
where stack and heap oftentimes grow in opposite directions due to dynamic
growth.

Segments allow to map each of the above blocks inpedently onto physical address
space at different locations, leading to better utilization of the physical
mmemory, reduced fragmentation, easy relocation to grow segments, or even
share segments between processes, e.g. code segment.

A segment is described by:

- *base*: the location of the segment in physical address space.

- *bounds*: size of the segment.

- *direction*: direction of growth (heap and stack grow in opposite directions).

- *offset*: offset of the segment in virtual address space. Code segment
  typically starts at 0B, then goes heap, and stack grows from the end of the
  virtual address space.

A virtual address is the memory location in the virtual address space. It also
includes the segment number to perform an address translation.

# EXAMPLE

```console
% segment -segA 2:2 -segB 4:2
configuration:
  virtual address bounds: 4KB
  SEG0 base: 2KB bounds: 2KB dir: positive virt-base: 0KB
  SEG1 base: 4KB bounds: 2KB dir: negative virt-base: 4KB

translations
 virt: 3462 (SEG1) phys: 5510
 virt: 3046 (SEG0) segmentation fault
 virt: 2484 (SEG0) segmentation fault
 virt: 2949 (SEG0) segmentation fault
 virt: 342 (SEG0) phys: 2390
```
