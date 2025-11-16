Name
====

**queue** - priority queue

Description
===========

Queue package holds different implementation of Priority Queue (PQ) with parameterized `less` function to create MaxPQ or MinPQ.

The discussions and comments imply MaxPQ.

Implementations
---------------

-	*Unordered array* [lazy]: store items in the array in undefined order. Push() appends a new item to the end. Search for the maximum item for Pop() and Top().

-	*Ordered array* [eager]: store items in a sorted array. Push() inserts the item at appropriate position (like insert sort). Pop() and Top() retrieve the last item.

-	*Binary heap* [eager]: store items in a heap-ordered complete binary tree. Push() appends the new item to the end and promotes it all the way to the root item as long is the parent is less than or equal to the new item. Pop() removes the root item, moves the last one to the root and demotes the new root with the max child. It continues demotion process until the end of array.

Performance
-----------

| Implementation  | Insert | Pop   |
|-----------------|--------|-------|
| Unordered array | 1      | N     |
| Ordered array   | N      | 1     |
| Binary heap     | ln(N)  | ln(N) |
