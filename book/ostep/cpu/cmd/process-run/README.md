<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**process-run** - schedule processes on a single CPU

Synopsis
========

```console
process-run -l spec
process-run -l spec -c
```

Description
===========

Process-run simulates scheduling of processes on a single CPU.

Arguments
---------

-	`-l spec`: a comma separated list of process specifications of the form `N:prob`, where `N` is the number of instructions, and `prob` is the probability of getting a CPU instruction else IO. The probability is an integer in the range[0, 100].

-	`-c`: run simulation.

Example
-------

Consider processes spec `3:50,5:100`. It creates two processes:

-	The first one has 3 instructions with 50% probability of getting a CPU instruction.
-	The second one has 5 CPU instructions.

Simulate run on the CPU:

```console
% process-run -l 3:50,5:100 -c
Process 0
  cpu
  io
  iod
  cpu
Process 1
  cpu
  cpu
  cpu
  cpu
  cpu
Clock   PID:0   PID:1
 1      run:cpu rdy
 2      blk:io  run:cpu
 3      blk:io  run:cpu
 4      blk:io  run:cpu
 5      rdy:iod run:cpu
 6      rdy:iod run:cpu
 7      run:cpu ok
 8      ok      ok
```

Note that the first process gets CPU,IO,CPU instructions. The `iod` means IO finished and the process is in ready state, waiting for the scheduler to pick it up.

The output sample reads:

-	The first process runs for one cycle (clk 1) while the second process waits in the ready state.
-	The first process blocks on the IO (clk 2) and remains in that state for the duration of IO (the cost is 3 cycles), while the second PID executes CPU instructions (clk 2-6).
-	The first process is unblocked on the IO (clk 5) and remains in the ready state (clk 5-6) until the second process finishes.
-	The second process completes and the remaining first process's CPU instruction runs (clk 7)
-	All processes complete at the end (clk 8)
