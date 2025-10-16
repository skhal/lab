# NAME

**lru** - least recently used cache


# DESCRIPTION

## Problem

Implement a Least Recently Used (LRU) cache of a given capacity C to store the
most recent C items. An item becomes most recent in one of the two cases:

  * when it is added to the cache
  * when it is accessed

API

```
type Cache interface {
  func Put(key, value int)
  func Get(key int) (value int, ok bool)
}
```

Create a new cache with positive integer capacity.

## Example

```
cache, _ := lru.NewCache(3)
cache.Put(1, 100)
cache.Put(2, 200)
cache.Put(3, 300)
cache.Put(4, 400)
```

Evicts least recent item with key 1. The content keys are in order (2, 3, 4)
with the least recent at the index 0 and the most recent at the index 2.
