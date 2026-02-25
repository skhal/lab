<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**scheduler** - schedule jobs with different policies

# DESCRIPTION

`scheduler` demonstrates different scheduling policies to run jobs.

## Terminology

Let's assume that a job has following time event:

- *arrival*: time when the job appears in the system, i.e., a process fork.
- *firstrun*: time when the scheduler picks up the job to run for the first
  time.
- *completion*: time when the job completes.

We can define following metrics:

```
Tturnaround = Tcompletion - Tarrival
```

Auxiliary:

- *workload*: a collection of outstanding jobs at a given moment in time

## Algorithms

**FIFO**: (First In First Out), schedules jobs to run in the order of arrival.
The selected job runs uninterrupted.

```console
% scheduler -job-spec 1,4,7
jobs: 3
scheduler: fifo

jobs:
  1 duration: 1
  2 duration: 4
  3 duration: 7

stats:
  1  Response:  0 Turnaround  1 Wait  0
  2  Response:  1 Turnaround  5 Wait  1
  3  Response:  5 Turnaround 12 Wait  5

average:
     Response:  2 Turnaround  6 Wait  2
```
