// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

// InstructionsSet holds instructions for the virtual machine. It is pretty
// much a post-order AST.
type InstructionsSet struct {
	Instructions []any // post-ordered AST instructions
}

// Number is a floating precision number.
type Number float64

// BinOp enumerates supported binary operators.
//
//go:generate stringer -type=BinOp -linecomment
type BinOp int

const (
	_ BinOp = iota
	// keep-sorted start
	BinOpMinus // minus
	BinOpPlus  // plus
	// keep-sorted end
)

// Ref is a cell reference.
type Ref string

// Function enumerates supported functions.
//
//go:generate stringer -type=Function -linecomment
type Function int

const (
	_ Function = iota
	// keep-sorted start
	FuncSum // SUM
	// keep-sorted end
)

var calls = map[string]Function{
	"SUM": FuncSum,
}

// Call is a function call, e.g. "SUM(...)".
type Call struct {
	Func Function // function identifier
	Args int      // the number of arguments
}
