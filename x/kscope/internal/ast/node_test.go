// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"fmt"

	"github.com/skhal/lab/x/kscope/internal/ast"
)

func ExampleNumber_String() {
	var n ast.Node = ast.Number{Val: 1.23}
	fmt.Println(n)
	// Output:
	// 1.2
}

func ExampleBinExpr_String() {
	var n ast.Node = ast.BinExpr{
		Op:    ast.BinOpPlus,
		Left:  ast.Number{Val: 1},
		Right: ast.Number{Val: 2},
	}
	fmt.Println(n)
	// Output:
	// 1.0 + 2.0
}

func ExampleCall_String() {
	var n ast.Node = ast.Call{
		Name: "demo",
		Args: []ast.Node{
			ast.Number{Val: 1},
			ast.Number{Val: 2},
		},
	}
	fmt.Println(n)
	// Output:
	// demo(1.0, 2.0)
}

func ExampleFunc_String() {
	var n ast.Node = ast.Func{
		Name: "demo",
		Body: []ast.Node{
			ast.Number{Val: 1},
		},
	}
	fmt.Println(n)
	// Output:
	// def demo()
	//   1.0
}
