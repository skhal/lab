# NAME

**flatten solution** - flatten a singly linked list with child nodes


# DESCRIPTION

> [!NOTE]
> This solution does not modify the input multi-layered list

Use queue to process the tree in layers. Start with the root node and enqueue
every child node.

Algorithm:

  - start with a queue with current node only
  - for each item in the queue:
    + generate a flat sub-list of values
    + generate a sub-queue of child nodes
	  + connect the tail of the main list to the head of the sub-list and update
      the tail position
    + append the sub-queue of child nodes to the main queue

Complexity:

  - time: O(n) to scan through all elements
  - space: O(n) to store all elements in a list and the queue
