Name
====

**38** - remove adjacent duplicate pairs

Description
===========

Problem
-------

Recursively remove adjacent duplicate pair characters from a string.

Example
-------

-	*Input*: "abccbaad"
-	*Output*: "d"

Removes:

-	"cc" to "abbaad"
-	"bb" to "aaad"
-	"aa" to "adt

Solution
--------

<details>
<summary>Details</summary>

Use stack to keep track of "saved" characters. When reading i-th character, skip
it if the last item on the stack is the same - pop that last item from stack as
well.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
