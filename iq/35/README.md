<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**35** - next largest number

Description
===========

Problem
-------

Given a list of numbers {ni}, create a parallel list {ki} where ki is equal to the first item in the {ni} list that is larger than ni or -1 if such item does not exist.

Example
-------

-	*Input*: [1, 3, 2, 5]

-	*Output*: [3, 5, 5, -1]

Solution
--------

<details>
<summary>Details</summary>

Observe that for a given item ni in the input list, we need to set it as the
next largest value for a set of previous numbers. Use stack to keep track of
the positions of these numbers.

Algorithm:

- For every number ni in the input list:
  * Set is as the next largest item for all items on the stack that keeps track
    of prior items if ni is larger.
  * Add i-index to the stack
- Mark left-over items on the stack as -1, i.e., next largest number does not
  exist.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
