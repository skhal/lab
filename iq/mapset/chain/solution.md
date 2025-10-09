# NAME

**chain solution** - solution to the longest chain in a collection


# DESCRIPTION

A naive solution is to sort the collection and linearly scan for the chains.
Time complexity would be `O(N*logN)` due to sort, space complexity is `O(n)` to
store the copy of the collection.

There is a linear solution in time complexity.

## Solution

The idea is to use the hash set to efficiently identify the beginning and walk a
chain.

Algorithm:

  * Store items in a set
  * For every item in the set that does not have a previous item, i.e., `n - 1`
    calculate the length of the chain.
  * Keep the longest chain.

Complexity:

  * *Time*: O(n)
  * *Space*: O(n)
