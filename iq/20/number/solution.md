<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**happy number solution** - is number happy

DESCRIPTION
===========

Observe that the sequence of a happy number products either ends in 1 or continues indefinitely with eventual cycle. If these numbers are viewed as a singly linked list, then one can apply standard cycle detection mechanism to the sequence.

The trick is that there is no physical link. Instead, we can use dynamically calculated numbers for two imaginable pointers, slow and fast, to detect a cycle.

Complexity:

-	*Time*: O(log N)
-	*Space*: O(1)
