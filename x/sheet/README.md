<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**sheet** -- demo of a spreadsheet engine

# DESCRIPTION

Sheet demonstrates a minimal engine to support spreadsheets, a table of cells
with formulas:

```go
s.Set("A1", "1")
s.Set("B1", "=SUM(A1:A5, 7-6)")
s.Set("B3", "=IF(B1 > 3, 10, 20)")
s.Calculate()
s.VisitAll(func(id, val string, res float64) bool {
	fmt.Printf("%s %3.1f\t%s\n", id, res, val)
	return true
})
```

The expected output is:

```
A1 1.0  1
B1 2.0  =SUM(A1:A5, 7-6)
B3 20.0 =IF(B1 > 3, 10, 20)
```

## Engines

It supports two engines to drive cell parsing and fromula calculation:

- *AST*: parse formulas into Abstrast Syntax Tree (AST). Walk the three in
  postorder traversal.

- *VM*: parse formulas into AST, then compile it into instructions. Use a
  Virtual Machine (VM) to execute the instruction to get the value.
