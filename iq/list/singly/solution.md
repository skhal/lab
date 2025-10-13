# NAME

**solution** - solution to reverse a singly linked list


# DESCRIPTION

Algorithm:

- keep track of the previous and current nodes initialized to `nil` and the
  start of the list.
- while current node is not `nil`, cache the next node reference, link current
  node to the previous one, and set previous to current and current to the next.

Complexity:

- time: O(n)
- space: O(1)
