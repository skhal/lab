# NAME

**intersection solution** - intersection of two singly linked lists


# DESCRIPTION

Consider two lists L1 and L2 with shared tail T and prefixes P1 and P2:

```
L1 = P1 + T
L2 = P2 + T
```

The objective is to traverse the lists in parallel until the nodes
compare the same (using pointers or other equality technique).

The only way for parallel traversal is by equating length of both prefixes.
Consider following sequences:

```
LA = L1 + L2 = P1 + T + P2 + T
LB = L2 + L1 = P2 + T + P1 + T
```

Notice that the last section T is common in both LA and LB.

Algorithm:

  * traverse in parallel L1+L2 and L2+L1 until common node is found

Complexity:

  * time: O(n+m)
  * space: O(1)
