<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**cycle solution** - detect cycles in a singly linked list

DESCRIPTION
===========

Run two iterators through the list at different speeds, say one step and two steps per iteration, and check whether the iterators point at the same item in the list. If cycle is present, the two pointers will eventually match.

Make sure to shortcut the loop if one of the iterators is nil.

Complexity:

-	*Time*: O(n) - scan all items
-	*Space*: O(1) - two pointers
