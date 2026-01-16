<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**33** - weighted random number

Description
===========

Problem
-------

Consider an array of size N. Each item is an integer value representing a weight of the number, given by the index in the array.

Write a function to get a random number from the array according to the provided weights.

Example
-------

*Input*: [4, 3, 2]

*Output*: generate value 0 in 4/9 cases, 1 in 3/9, and 2 in 2/9 cases.

Solution
--------

<details>
<summary>Details</summary>

Use binary search on a transformed array, where numbers represent cumulative
distribution function, i.e., a value at index k is the sum of weights for
indices 0 to k:

```
weights = [1, 2, 3, 4, ...]
cdf = [1, 3, 6, 10, ...]
```

Use a random number generator to get a random integer between 0 and the sum
of all weights (excluded).

Complexity:

- _time_: O(n + log(n)) ~ O(n) to build CDF and run binary search
- _space_: O(n) to store CDF

Optimizations:

- If the code re-uses the same weights to generate a large number of random
  numbers M, i.e., M >> N, then one can build the CDF once and reuse it M times
  to drive the time complexity to O(log(N)) and get rid of O(M) factor.

</details>
