// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import "fmt"

// Token is the smallest unit of input.
type Token struct {
	Val  string    // token value
	pos  int       // token position in the text
	Kind TokenKind // token kind
}

// String prints the token.
func (t Token) String() string {
	return fmt.Sprintf("%s %q", t.Kind, t.Val)
}

// TokenKind identifies the kind of a token.
//
//go:generate stringer -type=TokenKind -linecomment
type TokenKind int8

const (
	_ TokenKind = iota

	// keep-sorted start
	TokAssign  // assign
	TokComma   // comma
	TokComment // comment
	TokDef     // definition
	TokDiv     // divide
	TokExtern  // extern
	TokIdent   // identifier
	TokLpar    // left-parenthesis
	TokMinus   // minus
	TokMul     // multiply
	TokNum     // number
	TokPlus    // plus
	TokRpar    // right-parenthesis
	TokVar     // variable
	// keep-sorted end
)
