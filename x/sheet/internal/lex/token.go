// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

// Token is a single item in the formula.
type Token struct {
	Type TokenType // token type
	Text string    // token value
	Err  error     // only set for TokenError
}

// TokenType enumerates different tokens.
//
//go:generate stringer -type TokenType -linecomment
type TokenType int

const (
	_ TokenType = iota
	// keep-sorted start
	TokenComma          // comma
	TokenDivide         // divide
	TokenEqual          // equal
	TokenError          // error
	TokenGreater        // greater
	TokenGreaterOrEqual // greaterorequal
	TokenIdent          // identifier
	TokenLess           // less
	TokenLessOrEqual    // lessorequal
	TokenLpar           // rpar
	TokenMinus          // minus
	TokenMultiply       // multiply
	TokenNotEqual       // notequal
	TokenNumber         // number
	TokenPlus           // plus
	TokenRange          // range
	TokenRpar           // lpar
	// keep-sorted end
)
