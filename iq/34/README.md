<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**34** - validate parenthesis expression

Description
===========

Problem
-------

Validate a string of a parenthesis expression. Valid parenthesis pairs are "()", "[]", and "{}".

Example
-------

Expression "([])" is valid but "{[]" is invalid because "{" is unbalanced.

Solution
--------

<details>
<summary>Details</summary>

Use stack to keep track of opening parenthesis and validate last item on the
stack when facing a closing parenthesis.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
