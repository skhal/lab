Name
====

**31** - matrix search

Description
===========

Problem
-------

Given a matrix n-by-m of sorted numbers, where the first item of the row is not less then the last item of the previous row, find whether a given number is present in the matrix.

Example
-------

The matrix is \[[1 2 3], [4, 5, 6]]. Value 4 is present and value 7 is not.

Solution
--------

<details>
<summary>Details</summary>

Use binary search algorithm with range between 0 and k that is equal to `n * m`.
Use modulo division to translate the index into row and column in the matrix.

Complexity:

- _time_: O(log(n * m))
- _space_: O(1)

</details>
