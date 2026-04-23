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
type Number struct {
	Val float64 // number value
}

// Operator enumerates supported binary operators.
//
//go:generate stringer -type=Operator -linecomment
type Operator int

const (
	_ Operator = iota
	// keep-sorted start
	OpMinus // minus
	OpPlus  // plus
	// keep-sorted end
)

// BinOp is a binary operator. It expects two operands.
type BinOp struct {
	Op Operator // the operator
}

// Ref is a cell reference.
type Ref struct {
	Cell string // cell identifier, e.g. "A1"
}

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
