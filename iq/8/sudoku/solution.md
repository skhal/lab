NAME
====

**sudoku solution** - solution to Sudoku validator

DESCRIPTION
===========

The problem has auxiliary complexity from coming up with an efficient API to set the board.

Let a box be a single place on the board with a number. It has a column, A to I, and a row, 1 to 9:

https://github.com/skhal/lab/blob/15c24247a85c9fbee9ddc28e98c7e4e4649f2dbc/iq/mapset/sudoku/solution.go#L5-L17

Using chess notation, we can give every column an identifier in the form of `<column><row>`, e.g. A1 or C3. Let the Sudoku board constructor take a list of set boxes:

https://github.com/skhal/lab/blob/15c24247a85c9fbee9ddc28e98c7e4e4649f2dbc/iq/mapset/sudoku/solution_example_test.go#L12-L15

Solution
--------

Use sets to keep track of seen numbers in a row, column and a block, i.e., one set per entity.

Assuming box enumeration runs row-by-row then column-by-column, i.e., A1 is 0, A2 =1, B1 = 9, etc., it is easy to show that:

```
col = boxID / 9
row = boxID % 9
```

To get a block, that is enumerated left-to-right in each row:

```
    A B C   D E F   G H I
1 | . . . | . . . | . . . |
2 | . 0 . | . 1 . | . 2 . |
3 | . . . | . . . | . . . |
    -----   -----   -----
4 | . . . | . . . | . . . |
5 | . 3 . | . 4 . | . 5 . |
6 | . . . | . . . | . . . |
    -----   -----   -----
7 | . . . | . . . | . . . |
8 | . 6 . | . 7 . | . 8 . |
9 | . . . | . . . | . . . |
    -----   -----   -----

block = 3 * (row / 3) + (col / 3)
```

Complexity
----------

*Time*: O(n^2) - scan all boxes.

*Space*: O(n^2) - keep track of seen elements in every block.
