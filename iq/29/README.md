Name
====

**29** - search in a rotated sorted array

Description
===========

Problem
-------

Given an array of sorted numbers that is rotated by k positions, and a number x, find an index of x in the array. Return -1 if x is not present in the array.

*Definition*: an array is rotated by k positions if the last k elements are put to the head of the array. For example: array [4, 5, 1, 2, 3] is the array [1, 2, 3, 4, 5] rotated by 2 positions.

Example
-------

Given an array [4, 5, 1, 2, 3], value 5 has index 1 and value 6 has index -1.

Solution
--------

<details>
<summary>Details</summary>

Use binary search algorithm with optional shift in the opposite direction upon
checking the value of the midpoint if the values in the range [left, right) are
rotated.

Complexity:

- _time_: O(log(n))
- _space_: O(1)

</details>
