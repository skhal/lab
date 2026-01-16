<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**stripzero solution** - solution to in place zero out matrix problem

DESCRIPTION
===========

Solution
--------

Use two hash sets to keep track of rows and columns to zero out. Scan the matrix to fill these sets in a first pass and then zero out marked rows and columns.

### Complexity

*Time*: O(m * n)

*Space*: O(m + n) to track rows and columns

Optimization
------------

Observe that if a row or a column to be zeroed out, it zeroes the very first row or column respectively. Use this fact to keep track of the rows and columns to zero out instead of using hash sets.

Notes:

-	Use a flag to indicate whether to zero out the first row or column itself

### Complexity

*Time*: O(m * n)

*Space*: O(1)
