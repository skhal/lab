<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**27** - find lower bound

Description
===========

Problem
-------

Given a sorted array nn, find lower bound for a number x. A lower bound is defined as the minimal value in nn that is not less than x.

Example
-------

The array is [1, 3, 4]. Lower bound for 2 is 3, and for 3 is 3.

Solution
--------

<details>
<summary>Details</summary>

Use binary search to lookup for x but return the upper bound. Take care of the
cases when there is no upper bound available if the value x is above the maximum
number in the array.

</details>

See Also
========

-	C++ [std::lower_bound](https://en.cppreference.com/w/cpp/algorithm/lower_bound.html)
