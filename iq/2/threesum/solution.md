SOLUTION
========

Hint
----

A triplet sum `x1 + x2 + x3 = C` where C is some constant is equivalent to the problem of finding a pair [x2, x3] that satisfies the condition `x2 + x3 = C - x1`.

Algorithm
---------

1.	Sort the array in descending order
2.	Walk the array thus fixing `x1`
	-	inward scan the rest of the array for a pair [x2, x3] that satisfies the condition `x2 + x3 = C - x1`

Complexity
----------

*Time*: it is roughly equivalent to sort and scan, e.g. O(n log(n)) and O(n^2). The second term wins leading to O(n^2).

*Space*: O(n) to hold the copy of the array for sorting.

OPTIMIZATIONS
=============

If `C = 0`, at least one value must have the opposite sign for the condition `x1 + x2 + x3 = 0` to hold.

1.	The beginning and the end of the array must of opposite signs.
2.	Stop the scan when x1 becomes positive.

More optimizations:

1.	Skip the same x1 items
2.	Skip the same x2 items
3.	Skip the same x3 items
