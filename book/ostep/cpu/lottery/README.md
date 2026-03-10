<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**lottery** - lottery job scheduler

# DESCRIPTION

A lottery scheduler randomly picks up next job to run using weights assigned to
every job.

It is called lottery because this implementation uses *tickets*, an integer
number, to emulate weights. The number of assigned tickets is arbitrary, but
it defines process weight:

```
weight_i = tickets_i / sum(tickets_i)
```

## Example

```console
% lottery -jobs 10:10,10:100,10:200
jobs:
  jid:1 len:10 tks:10
  jid:2 len:10 tks:100
  jid:3 len:10 tks:200

trace:
  1     jid:2 len:10 tks:100
  2  +4 jid:3 len:10 tks:200
  7     jid:2 len:10 tks:100
  8  +2 jid:3 len:10 tks:200
  11    jid:2 len:10 tks:100
  12 +1 jid:3 len:10 tks:200 [done]
  14 +6 jid:2 len:10 tks:100 [done]
  21 +9 jid:1 len:10 tks:10 [done]
```
