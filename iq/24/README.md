<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**24** - find the insertion index

Description
===========

Problem
-------

Given a sorted array with unique values and a target number, return the index of the given number. If the number is not found, return the index where the number should be placed.

Examples
--------

The input array is [1, 3, 5, 7]. Consider following scenarios:

-	the target number is 3 and the return value is 1.
-	the target number is 4 and the return value is 2.

Solution
--------

<details>
<summary>Solution</summary>

Use binary search to lookup for the target number. If the number is found,
return the index, else return the next index of the last checked number.

Complexity:

- _time_: O(log(n))
- _space_: O(1)

</details>
