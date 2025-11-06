Name
====

**40** - max in a sliding window

Description
===========

Problem
-------

Consider an array of size N. Let window be a set of M consequent items from the array.

Find maximum value for a window of size M that slides through the array from left to right.

Example
-------

*Input*: [1, 4, 2, 3, 1, 2], window size is 3

*Output*: [1, 4, 4, 4, 3, 3]

Solution
--------

<details>
<summary>Details</summary>

Consider a window of size M with items at indices [i, i+M). Assume that the
window stores the values in descending order in a container C, with the front
one being the maximum value.

Let's slide the window by one position, to cover items at positions [i+1, i+1+M).
The goal is to rebalance sorted values in C.

(A) Consider old item at position i, that is removed from the window:

1. If the value is the maximum, remove it from C.

(B) Consider the new item at position i+M:

1. None of the items in C with values less or equal to it contribute to the next
   M windows - remove these items from C.

2. The container C has at most values that are greater than the new one. Push it
   to the end of the container C - it will be a new candidate in the future.

The logic works great if the input array holds unique values. It breaks in A.1 if
the maximum value is repeated. We want to keep track of the index of the items
in order to remove the max value only if it has the index i.

Performance:

- _time_: O(n)
- _space_: O(m) where m is the size of the window

</details>
