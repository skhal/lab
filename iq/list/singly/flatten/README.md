# NAME

**flatten** - flatten a singly linked list with child nodes


# DESCRIPTION

## Problem

Consider a singly linked list where a node may also have a single child node
that is the same type singly linked list with own possible child node. That
is a node has:

  - a value
  - a reference to the next node
  - a reference to the child node

Write a program that flattens such multi-level list into a single level. It
should process multi-level list by layers.

## Example

Input:

```
L1  L2  L3
1
2 - 6 - 9
    7
3
4 - 8 - 10
        11
5
```

Output:

```
[1 2 3 4 5 6 7 8 9 10 11]
```


# SEE ALSO

* [Solution](./solution.md)
