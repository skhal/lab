Name
====

**10** - find longest chain in a collection

Description
===========

Problem
-------

Given an array of numbers, find the longest chain of consecutive numbers. Two numbers are consecutive if they have a difference of 1.

Example
-------

**Input**: [7, 1, 8, 9, 2, 12]

**Output**: [7, 8, 9]

Solution
--------

<details>
<summary>Details</summary>

A naive solution is to sort the collection and linearly scan for the chains.
Time complexity would be `O(N*logN)` due to sort, space complexity is `O(n)` to
store the copy of the collection.

There is a linear solution in time complexity.

The idea is to use a hash set to efficiently identify the beginning and walk
the chain.

Algorithm:

- Store items in a set
- For every item in the set that does not have a previous item, i.e., `n - 1`
  calculate the length of the chain.
- Keep the longest chain.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
