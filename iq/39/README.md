Name
====

**39** - queue using stack

Description
===========

Problem
-------

Implement a queue using stack with the following API, derived from [`std::queue`](https://en.cppreference.com/w/cpp/container/queue.html):

```
type Queue interface {
    Empty() bool
    Front() (int, bool)
    Pop()
    Push(int)
    Size() int
}
```

Example
-------

Solution
--------

<details>
<summary>Details</summary.

Use two stacks:

- _in_: keeps track of input items. It gives a way to quickly push items to the
  queue.

- _out_: a list of items to remove - the top item is the first one in the queue.

Performance:

- _time_:
  * Empty is O(1)
  * Front is O(n) due to move from in to out stacks
  * Pop is O(n) due to move from in to out stacks
  * Push is O(1)
  * Size is O(1)

- _space_: O(n)

</details>
