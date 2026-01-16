<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

SOLUTION
========

Algorithm
---------

Initialization:

-	Find the position of the first zero value item, `zi`.
-	For each non-zero value position in the remaining array `ni`:
	-	Swap the items at positions `zi` and `ni`
	-	Advance `zi` by one

Complexity
----------

*Time*: linear scan, O(n)

*Space*: constant number of variables per iteration, O(1)
