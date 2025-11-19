Name
====

**singly** - singly linked list

Problem
=======

Reverse a singly linked list.

Example
-------

*Input*: [1, 2, 3]

*Output*: [3, 2, 1]

Solution
========

<details>
<summary>Details</summary>

Algorithm:

-   keep track of the previous and current nodes initialized to `nil` and the start of the list.
-   while current node is not `nil`, cache the next node reference, link current node to the previous one, and set previous to current and current to the next.

Complexity:

-   _Time_: O(n)
-   _Space_: O(1)

</details>
