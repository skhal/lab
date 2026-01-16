<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**23** - find longest uniform substring after k substitutes

Description
===========

Problem
-------

*Definition*: A uniform string consists of the same characters.

Given a string s and the number n, find longest uniform string that results after at most n substitutions.

Example
-------

Input string is "abacde" and n is 2. The longest string is "abac".

Solution
--------

<details>
<summary>Details</summary>

Use sliding window:

- Keep track of frequencies of characters in the window when
growing and shrinking it.
- Track maximum frequency in the window.
- The number of replacements is equal to the window size minus the maximum
    frequency.

Slide the window when the number of replacements surpasses the number of allowed
substitutes.
There is no need to update the maximum frequency since the algorithm is looking
for the longest window. The next larger window would trigger when slide and get
valid number of substitutes because a new letter becomes most frequent and
automatically drops the number of substitutions.

Complexity:

- _Time_: O(n)
- _Space_: O(n)

</details>
