Name
====

**28** - find upper bound in a sorted array

Description
===========

Problem
-------

Given an array of sorted integer numbers, find upper bound for a given number x.

*Definition*: an upper bound is a maximum value in the array that is less than a given value.

Example
-------

Input array is [1, 3, 5]. The upper bounds are:

-	1 for 2 and 3
-	3 for 4 and 5
-	5 for 6

Solution
--------

<details>
<summary>Details</summary>

Use binary search to lookup for x but return the lower bound. Take care of the
cases when there is no lower bound available for x below the lowest number in
the array.

</details>

See Also
========

-	C++ [std::upper_bound](https://en.cppreference.com/w/cpp/algorithm/upper_bound.html)
