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

Performance
-----------

| Implementation  | Insert | Pop |
|-----------------|--------|-----|
| Unordered array | 1      | N   |
| Ordered array   | N      | 1   |
