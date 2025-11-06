Name
====

**37** - remove adjacent duplicates

Description
===========

Problem
-------

Recursively remove adjacent duplicate characters from a string.

Example
-------

-	*Input*: "abccbad"
-	*Output*: "d"

Removes:

-	"cc" to "abbad"
-	"bb" to "aad"
-	"aa" to "d"

Solution
--------

<details>
<summary>Details</summary>

Use stack to keep track of "saved" characters. When reading i-th character, make
sure to skip it and any looking-forward adjacent duplicates (to be skipped) and
looking-back adjacent duplicates (to be removed from the stack).

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
