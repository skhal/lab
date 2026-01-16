<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**two sum solution** - solution to pair sum in unsorted array

DESCRIPTION
===========

Brute force
-----------

There are two scans:

1.	Scan the array to build a map of values to indices.
2.	Scan the array a second time to find a complement number for a given item, i.e., for a given item `n1` find an item `target - n1`. If no such item exists, scan for the next item in the array.

Optimal
-------

Notice that if `n1` is given, then we need to efficiently find `n2 = target - n1`, and vice versa - if `n2` is given, the problem reduces to finding `n1 = target - n2`.

This observation has two corollaries:

-	Given `n1` we are looking ahead for `n2`
-	Given `n2` we are looking back for `n1`

The brute-force solution uses the first fact, therefore pre-build a map to keep track of all items-to-index translation to lookup for `n2` for a given `n1`.

The optimal solution makes a single pass and looks for `n1` for a given `n2` in a hash map.

Complexity
----------

*Time*: O(n)

*Space*: O(n) to store all items in a hash map in the worst case.

> [!NOTE] The optimal solution will short-cut more often on average than the brute-force solution.
