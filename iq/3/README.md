Name
====

**palindrome** - is a palindrome

Problem
=======

Given a string, detect whether it is a palindrome, i.e. reads same left-to-right and right-to-left.

Assumptions:

-	The input may include non-alphanumeric characters, not used in the check.
-	Treat upper and lower case the same.

Example:

-	*Input*: "Don't nod"
-	*Output*: Yes, palindrome

Solution
========

<details>
<summary>Details</summary>

Algorithm:

1.  Scan for the left and right rune skipping non-alphanum runes
2.  Compare left and right runes after converting to lower case
3.  Continue 1-2 as long as left position is below the right one

Complexity:

- _Time_: linear scan, O(n)
- _Space_: constant number of variables per iteration, O(1)

</details>
