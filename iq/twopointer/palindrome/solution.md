# SOLUTION

## Algorithm

1. Scan for the left and right rune skipping non-alphanum runes
2. Compare left and right runes after converting to lower case
3. Continue 1-2 as long as left position is below the right one

## Complexity

*Time*: linear scan, O(n)

*Space*: constant number of variables per iteration, O(1)
