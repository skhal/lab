<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**sched** - schedule jobs with different policies

# DESCRIPTION

`sched` demonstrates different scheduling policies to run jobs.

## Metrics

Consider the following job events:

- *arrival*: a moment in time when the job appears in the system.
- *firstrun*: the scheduler policy picks up the job to run for the first
  time.
- *completion*: the job completes.

Define following metrics:

```
Tresponse   = Tfirstrun - Tarrival
Tturnaround = Tcompletion - Tarrival
```

It is also helpful to track time the job spends idle, not running:

```
Twait = Tturnaround - Trunning
```

where `Trunning` is the collective time the job spends running on the CPU.

## First in First out policy (FIFO)

Run jobs in the order of arrival, uninterrupted.

```console
% go run ./book/ostep/cpu/cmd/sched/ -job-spec 7,2:4,3:1,5:1 -policy fifo -trace
policy: fifo

jobs:
  1 arrival: 0 duration: 7
  2 arrival: 2 duration: 4
  3 arrival: 3 duration: 1
  4 arrival: 5 duration: 1

trace:
  0 run 1 for 7 cycles
  7 run 2 for 4 cycles
  11 run 3 for 1 cycle [done]
  12 run 4 for 1 cycle [done]

stats:
  1  Response: 0   Turnaround: 7   Wait: 0   [done]
  2  Response: 5   Turnaround: 9   Wait: 5   [done]
  3  Response: 8   Turnaround: 9   Wait: 8   [done]
  4  Response: 7   Turnaround: 8   Wait: 7   [done]

average:
     Response: 5   Turnaround: 8   Wait: 5
```

## Shortest Job First policy (SJF)

Out of a collection of pending jobs, run the job with shortest time to
complete, uninterrupted.

```console
% go run ./book/ostep/cpu/cmd/sched/ -job-spec 7,2:4,3:1,5:1 -policy sjf -trace
policy: sjf

jobs:
  1 arrival: 0 duration: 7
  2 arrival: 2 duration: 4
  3 arrival: 3 duration: 1
  4 arrival: 5 duration: 1

trace:
  0 run 1 for 7 cycles
  7 run 3 for 1 cycle [done]
  8 run 4 for 1 cycle [done]
  9 run 2 for 4 cycles

stats:
  1  Response: 0   Turnaround: 7   Wait: 0   [done]
  2  Response: 7   Turnaround: 11  Wait: 7   [done]
  3  Response: 4   Turnaround: 5   Wait: 4   [done]
  4  Response: 3   Turnaround: 4   Wait: 3   [done]

average:
     Response: 3   Turnaround: 6   Wait: 3
```

Notice how all average metrics improve in SJF compared to FIFO. Running short
jobs first improves system responsiveness.

## Shortest Time to Complete First policy (STCF)

Preempt running job to pick up new shortest time to complete job or continue
running current job if it is already the shortest one.

```console
% go run ./book/ostep/cpu/cmd/sched/ -job-spec 7,2:4,3:1,5:1 -policy stcf -trace
policy: stcf

jobs:
  1 arrival: 0 duration: 7
  2 arrival: 2 duration: 4
  3 arrival: 3 duration: 1
  4 arrival: 5 duration: 1

trace:
  0 run 1 for 2 cycles
  2 run 2 for 1 cycle
  3 run 3 for 1 cycle [done]
  4 run 2 for 1 cycle
  5 run 4 for 1 cycle [done]
  6 run 2 for 2 cycles
  8 run 1 for 5 cycles

stats:
  1  Response: 0   Turnaround: 13  Wait: 6   [done]
  2  Response: 0   Turnaround: 6   Wait: 2   [done]
  3  Response: 0   Turnaround: 1   Wait: 0   [done]
  4  Response: 0   Turnaround: 1   Wait: 0   [done]

average:
     Response: 0   Turnaround: 5   Wait: 2
```

This policy strives to reduce idle time by bringing the best of SJF via job
preemption mechanism.

## Round Robin policy (RR)

Give every pending job a chance to run for a single cycle by constantly
switching between the jobs.

```console
% go run ./book/ostep/cpu/cmd/sched/ -job-spec 7,2:4,3:1,5:1 -policy rr -trace
policy: rr

jobs:
  1 arrival: 0 duration: 7
  2 arrival: 2 duration: 4
  3 arrival: 3 duration: 1
  4 arrival: 5 duration: 1

trace:
  0 run 1 for 2 cycles
  2 run 2 for 1 cycle
  3 run 1 for 1 cycle
  4 run 3 for 1 cycle [done]
  5 run 2 for 1 cycle
  6 run 1 for 1 cycle
  7 run 4 for 1 cycle [done]
  8 run 2 for 1 cycle
  9 run 1 for 1 cycle
  10 run 2 for 1 cycle [done]
  11 run 1 for 2 cycles

stats:
  1  Response: 0   Turnaround: 13  Wait: 6   [done]
  2  Response: 0   Turnaround: 9   Wait: 5   [done]
  3  Response: 1   Turnaround: 2   Wait: 1   [done]
  4  Response: 2   Turnaround: 3   Wait: 2   [done]

average:
     Response: 0   Turnaround: 6   Wait: 3
```
