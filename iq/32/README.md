<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**32** - find local maxima

Description
===========

Problem
-------

Given an array, find local maxima.

*Definition*: local maxima is a value that is greater than its immediate neighbors.

Assumptions:

-	the immediate neighbors are guaranteed to be non-equal numbers

Example
-------

*Input*: [4 5 1 3 2]

*Output*: valid answers are 5 or 3.

Solution
--------

<details>
<summary>Details</summary>

Use binary search. Check neighbors to detect the slope and move in corresponding
direction.

Complexity:

- _time_: O(log(n))
- _space_: O(1)

</details>
