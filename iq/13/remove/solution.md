<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**remove solution** - solution to remove elements from singly linked list

DESCRIPTION
===========

Algorithm
---------

The objective is to find previous node to the node to be removed. Keep in mind that it may not exist if the node to be removed is the head of the list.

One option is to create a fake node to point to the head of the list to guarantee existence of the previous node. Another option is to use a boolean flag to indicate whether the head should be removed.

Algorithm:

-	skip first N nodes
-	keep track of the previous node starting from N+1 nodes

Complexity:

-	*time*: O(n) to scan the list
-	*space*: O(1) store a single reference to the previous node
