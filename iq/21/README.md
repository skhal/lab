<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**21** - find anagrams in a string

Problem
=======

Given a strings s and t find all anagrams of t inside s.

Example
-------

String s is "abaab". String t is "ab".

Output is a list of strings ["ab", "ba", "ab"].

Solution
========

<details>
<summary>Details</summary>

Use a fixed sliding window to scan through the string to check whether sub-strings are anagrams of the target string.

Complexity, assuming s has length n and t has length k:

-   _Time_: O(n) - scan through characters in the string.
-   _Space_: O(1) - constant space for character frequency storage.

</details>
