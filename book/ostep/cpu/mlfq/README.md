<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**mlfq** - Multilevel Feedback Queue scheduling policy

# DESCRIPTION

Mlfq implements a Multilevel Feedback Queue scheduling policy
([ostep-mlfq](https://github.com/remzi-arpacidusseau/ostep-homework/tree/afb36ca8ddbf81d847d18f6bd18a87f0a18667f2/cpu-sched-mlfq))
with following rules:

- add new processes to the top-priority queue (slice)
- round robin processes starting from the highest priority queue
- if a process uses allotted CPU cycles at priority P, deprioritize it to
  lower priority queue at P+1
- move all processes to the top-priority every N cycles

The process may trigger an IO every once in a while (see `proc.Spec.IOCycles`).
The policy skips such processes until the IO is done.

```console
% go run ./book/ostep/cpu/mlfq/ -policy 1:5:20 -proc 0:10,2:10:3:2,5:5:2:10
policy:
  allotment   1
  priorities  5
  boost       20

processes:
  1 arrive:0 cpu:10
  2 arrive:2 cpu:10 io:2 after:3
  3 arrive:5 cpu:5 io:10 after:2

trace:
  1 run 1 cpu 1 pri 0
  2 run 1 cpu 2 pri 1
  3 run 2 cpu 1 pri 0
  4 run 2 cpu 2 pri 1
  5 run 1 cpu 3 pri 2
  6 run 3 cpu 1 pri 0
  7 run 3 cpu 2 pri 1 [block]
  8 run 2 cpu 3 pri 2 [block]
  9 run 1 cpu 4 pri 3
  10 run 1 cpu 5 pri 4
  11 run 2 cpu 4 pri 3
  12 run 2 cpu 5 pri 4
  13 run 1 cpu 6 pri 4
  14 run 2 cpu 6 pri 4 [block]
  15 run 1 cpu 7 pri 4
  16 run 1 cpu 8 pri 4
  17 run 2 cpu 7 pri 4
  18 run 3 cpu 3 pri 2
  19 run 3 cpu 4 pri 3 [block]
  20 run 2 cpu 8 pri 0
  21 run 1 cpu 9 pri 0
  22 run 2 cpu 9 pri 1 [block]
  23 run 1 cpu 10 pri 1 [done]
  24 -
  25 run 2 cpu 10 pri 2 [done]
  26 -
  27 -
  28 -
  29 -
  30 run 3 cpu 5 pri 0 [done]

stats:
  1 response: 1 turnaround: 23 wait: 13
  2 response: 1 turnaround: 23 wait: 13
  3 response: 1 turnaround: 25 wait: 20

average:
    response: 1 turnaround: 23 wait: 15
```

## Notes on implementation

Unlike
[ostep-mlfq](https://github.com/remzi-arpacidusseau/ostep-homework/tree/afb36ca8ddbf81d847d18f6bd18a87f0a18667f2/cpu-sched-mlfq)
this implementation is broken into packages for readability and modularity.
It helps to decouple concepts, break the code into small, testable parts, and
put Go best practices to use. It really helps to better understand coupling
between components and reveal the structure of the simulator.

Packages:

- `internal/cmd`: implements the command - parse command line flags, run
  simulator, and generate the report.

- `internal/sim`: drives the simulation. For For each cycle C, it completes the
  previous cycle C-1: unblock IO blocked processes and schedules new processes.
  Then it advances to the next cycle to run CPU and IO. This is where the
  simulator integrates with the policy to get the next process to run and
  accounts for the IO.

- `internal/policy`: MLFQ scheduler groups processes into slices by priority,
  starting with the top-priority (0) (aka priority queues in this problem).
  MLFQ round robins processes in each priority to give a chance to run every
  process in the queue. Once the process uses allotted CPU cycles, the policy
  lower the processes's priority.

  It also accounts for blocked processes here by skipping such processes when
  scanning through the priority queue.

- `interface/queue`: implemnets a round-robin queue using slices.

- `interface/io`: emulates an IO to account for iO cycles.

- `interface/cpu`: gives access to CPU clock (cycles).

- `interface/proc`: emulates process. There are two interfaces:
  `process.Process` and `proc.Control`. The former one is read-only view at
  the process, designed for the reporter and other components that only need
  to access process's state. The latter one is for the simulator and other
  components that must control the process, i.e., mutate the process's state.
