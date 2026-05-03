// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import "fmt"

// Token is the smallest unit of input.
type Token struct {
	Val  string    // token value
	Pos  Position  // token position in the text
	Kind TokenKind // token kind
}

// String prints the token.
func (t Token) String() string {
	return fmt.Sprintf("%s: %s %s", t.Pos, t.Kind, t.Val)
}

// TokenKind identifies the kind of a token.
//
//go:generate stringer -type=TokenKind -linecomment
type TokenKind int8

const (
	_ TokenKind = iota

	// keep-sorted start
	TokComm  // comment
	TokDef   // definition
	TokDiv   // divide
	TokExt   // extern
	TokIdent // identifier
	TokMinus // minus
	TokMul   // multiply
	TokNum   // number
	TokPlus  // plus
	// keep-sorted end
)

// Position identifies the location of token in the input string. It holds the
// index where the token starts and the index of one past last character of the
// token.
//
// For example, token "123" in string "abc 123 def" has position {4, 7}.
type Position struct {
	Start int // location of the token
	End   int // location of the past last symbol of the token
}

// String prints the position.
func (p Position) String() string {
	return fmt.Sprintf("%d..%d", p.Start, p.End)
}
