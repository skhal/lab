NAME
====

**palindrome solution** - is singly linked list a palindrome

DESCRIPTION
===========

Solution
--------

There are two ways to compare nodes: move inward from both ends or move in opposite directions from the middle point.

This solution focuses on the latter approach, move in opposite directions. It also create auxiliary nodes to track backward direction to avoid mutations of the input list - these nodes link to each other and store a reference to the wrapped node from the input list.

Middle point algorithm:

-	Assume a list of N elements.
-	Linearly scan through elements with two references: r1 advances with every node, r2 advances with every second node
-	In the end, the r1 points at the tail of the list (it scans N elements) whereas r2 points at the middle of the list (it scans N/2 elements)

The solution uses the middle point algorithm to also build the auxiliary linked list for the left-part of the list.

Complexity
----------

-	time: O(n) two linear scans to find middle point and compare values.
-	space: O(n) for auxiliary list
