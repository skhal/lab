NAME
====

**8** - verify sudoku board

DESCRIPTION
===========

Problem
-------

Validate the current stat of a partially solved Sudoku board. The rules of the game are:

-	Every row and column must include unique numbers.
-	The grid is divided into 3x3 blocks.
-	Each 3x3 block must include unique numbers.

Example
-------

**Input**: consider following board, where dot indicates empty space

```
    A B C   D E F   G H I
    -----   -----   -----
1 | 1 . . | 2 . . | 3 . .
2 | . . . | . . . | . . .
3 | . . . | . . . | . . .
    -----   -----   -----
4 | . . . | . . . | . . .
5 | . . . | 2 . . | . . .
5 | . . . | . . . | . . .
    -----   -----   ----
7 | . . 4 | . . . | 4 . .
8 | . . . | . . 3 | . . .
9 | . . . | . 3 . | . . .
```

**Output**: `False` because either the column D holds 2 twice, or row 7 holds 4 twice, or block D7 holds 3 twice.

SEE ALSO
========

-	[Solution](./solution.md)
