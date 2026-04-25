// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

import "encoding/gob"

func init() {
	gob.Register(InstructionsSet{})
	gob.Register(InstType(0))
	gob.Register(Inst{})
	gob.Register(BinOp(0))
	gob.Register(Function(0))
	gob.Register(Call{})
}

// InstructionsSet holds instructions for the virtual machine. It is pretty
// much a post-order AST.
type InstructionsSet struct {
	Instructions []Inst // post-ordered AST instructions
}

// InstType defines the instruction type.
//
//go:generate stringer -type=InstType -linecomment
type InstType int

const (
	_ InstType = iota
	// keep-sorted start
	InstTypeBinOp  // operator
	InstTypeCall   // function
	InstTypeNumber // number
	InstTypeRef    // reference
	// keep-sorted end
)

// Inst is an instruction. It acts like a union: only one field is set
// depending on the instruction type.
type Inst struct {
	Type   InstType // instruction type
	Number float64  // number literal
	BinOp           // binary operator
	Ref    string   // cell reference
	*Call           // function call
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
	BinOpDivide   // dividie
	BinOpMinus    // minus
	BinOpMultiply // multiply
	BinOpPlus     // plus
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
