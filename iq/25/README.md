Name
====

**25** - range of indices

Description
===========

Problem
-------

Given an array of sorted numbers with duplicates and a given number n, find a range of indices where the number n is present. Return (-1, -1) if n is not present in the array.

Example
-------

Consider array [1, 2, 2, 3] and a number 2. The range would be (1, 2).

Solution
--------

<details>
<summary>Solution</summary>

Use binary search to find lower and upper bounds of the number.

Complexity:

- _time_: O(log(n))
- _space_: O(1)

</details>
