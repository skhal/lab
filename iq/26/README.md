Name
====

**26** - max cut of number

Description
===========

Problem
-------

Consider an array of integer numbers {ni} and a number k. Find the maximum integer value m such that for all ni > m the following condition holds:

```
sum(ni) >= k
```

Example
-------

The array is [1, 3, 2, 4] and k = 2. The maximum value is 2 because the sum of higher numbers above 2 is:

```
(3 - 2) + (4 - 2) = 1 + 2 = 3
```

Solution
--------

<details>
<summary>Details</summary>

Let's create a function `sumOver(l)` to calculate `sum(ni)` for `ni > l`, where
`l` is a some number between 0 and the maximum number in the array.

The problem can be rephrased as: find maximum number `nu` between 0 and the
maximum value in the array such that `sumOver(nu)` >= k. That is use binary
search to search for `nu` and `sumOver(nu)` to decide to move left or right.

Complexity:

- _time_: O(n log(m)) where m is the maximum value in the array.
- _space_: O(1)

</details>
