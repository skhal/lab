<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

Name
====

**11** - find triplets of geometric sequences in a collection

Problem
=======

Given a collection of numbers, find the number of triples of [geometric sequences](https://en.wikipedia.org/wiki/Geometric_progression) in a collection such that items at indices i, j, and k satisfy the following condition: i < j < k.

Example
-------

*Input*: [3, 1, 2, 3, 9, 3, 27] and ratio is 3

*Output*: 3 because of triples [0, 4, 6], [3, 4, 6], and [1, 3, 4]

Solution
========

<details>
<summary>Details</summary>

For a given ratio `r` and number `n`, a geometric triplet is a sequence of numbers `[n, r * n, (r * r) * n]`.

We need to find a way to efficiently lookup for two numbers from the geometric triplet, once the third one is fixed.

If the first item is fixed `n`, then a hash map can help us find the positions of `r * n` and `(r * r) * n`. There is no way to guarantee that found indices satisfy condition `i < j < k`, assuming i, j, and k are fore the first, second, and the third item in the triplet. The algorithm would need to run additional checks to validate the indices.

The same performance penalty hit applies to fixing the third item in the triplet.

However, if we fix the second element in the triplet, `r * n`, and use two separate hash maps for the previous and next items in the input collection, we are guaranteed that `i < j < k` is satisfied.

The total number of triples for a given `r * n` item is the number of permutations between the previous set of `n` and next set of `(r * r) * n`. Assuming there are `m_prev` and `m_next` items, total number of permutations is:

```
m_prev * m_next
```

Algorithm:

-   Linearly scan the collection.
-   Keep track of the items on the left and right in separate hash maps - the key is the number, the value is the number of times it is present in the sub-set.
-   When moving to the next item in the collection, move the element between the tho hash maps.

Complexity:

- _Time_: O(n) linear scan.
- _Space_: O(n) two hash maps for left and right sub-collections.

</details>
