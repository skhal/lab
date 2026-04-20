// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

// Token is a single item in the formula.
type Token struct {
	Type tokenType // token type
	Text string    // token value
	Err  error     // only set for TokenError
}

//go:generate stringer -type tokenType -linecomment
type tokenType int

const (
	_           tokenType = iota
	TokenError            // error
	TokenNumber           // number
	TokenPlus             // plus
	TokenMinus            // minus
)
