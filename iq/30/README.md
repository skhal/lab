<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**30** - median of two sorted arrays

Description
===========

Problem
-------

Given two sorted arrays, find median value of the two.

Example
-------

Arrays [1 4 6] and [2 7 8], median is 5.

Solution
--------

<details>
<summary>Details</summary>

The problem is about finding ways to split the two arrays. When one array is
split into two halves, the second one is automatically split too because median
is guaranteed to split all numbers into two parts, 50% of numbers below and
above the number.

The binary search algorithm should check the split logic to detect a move to
the left, right, or stop. Since we have two splits, one per array, a correct
split should hold conditions that left side of the split is lower than the
right side within an array, and cross arrays.

Be careful when calculating the median to consider the case of even and odd
total number of items in the arrays combined.

</details>
