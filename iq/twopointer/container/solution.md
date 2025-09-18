# SOLUTION

## Algorithm

1. Start with end-to-end container
2. Advance pointers
  * left pointer if left height is smaller
  * or right pointer if right height is smaller
  * else both
3. Store new volume is larger
4. repeat 2-3

## Complexity

*Time*: linear scan - O(n)

*Space*: a constant number of variables per iteration - O(1)
