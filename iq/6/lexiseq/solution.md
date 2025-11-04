SOLUTION
========

Algorithm
---------

1.	Find pivot scanning the string from right. The pivot is the first character at position `i` that satisfies the following condition: `s[i] < s[i+1]`. All characters to the right from the pivot are in descending order.
2.	If no pivot is found, reverse the string - end.
3.	Find the first character at `j`, when scanned from the right, that satisfies the following condition: `s[i] < s[j]`.
4.	Swap i and j characters.
5.	Reverse the right side from the pivot.

Complexity
----------

*Time*: 2 linear scans if no pivot is found, else 3 linear scans, O(n)

*Space*: a buffer to store the result, O(n)
