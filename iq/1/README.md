Name
====

**1** - pair sum in a sorted array

Description
===========

Problem
-------

Given an array of sorted numbers find indices of a pair of numbers that adds up to a given value `S`.

EXAMPLE
-------

-	Array is [1, 2, 3, 4, 5]
-	S is 5
-	Output is [0, 3]

Solution
========

<details>
<summary>Details</summary>

Apply inward scan with two indices starting from the opposite directions. Given that the array is sorted, say in ascending order, move one or the other index depending on the current sum value compared to the target value `S`.

Complexity:

- _time_: O(n)
- _space_: O(1)

</details>
