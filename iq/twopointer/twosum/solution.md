# SOLUTION

Apply inward scan with two indices starting from the opposite directions. Given
that the array is sorted, say in ascending order, move one or the other index
depending on the current sum value compared to the target value `S`.

## Complexity

*Time*: a linear scan is O(n).

*Space*: constant factor of variables per iteration is O(1).
