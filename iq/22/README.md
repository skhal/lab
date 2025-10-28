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
<summary>Simple</summary>

Use a sliding window with left and right sides. Keep track of seen characters
in a hash map.

When expanding the window, move the right edge as long as there are new
characters. Store the duplicate character, call it stop-character. To shrink the
window, advance the left edge until the stop-character while cleaning seen-map.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>

<details>
<summary>Optimized</summary>

This solution uses the same idea as the simple solution but it changes how the
seen-map works. Instead of keeping track of seen characters in the set, it
stores the positions of seen characters.

It becomes easy to detect whether a character is within the window: if the index
of the newly read character is within the window bounds, then it is a duplicate
character, otherwise it is outside the window and can be overwritten with the
new position within the window.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
