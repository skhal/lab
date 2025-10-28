Name
====

**22** - find longest substring with unique characters

Description
===========

Problem
-------

Given a string, find longest substring that consists of unique characters.

Example
-------

Input string is "abcad". The longest sub-string is "bcad".

Solution
--------

<details>
<summary>Solution</summary>

Use a sliding window with left and right sides. Keep track of seen characters
in a hash map.

When expanding the window, move the right edge as long as there are new
characters. Store the duplicate character, call it stop-character. To shrink the
window, advance the left edge until the stop-character while cleaning seen-map.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
