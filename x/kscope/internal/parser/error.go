// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"errors"

	"github.com/skhal/lab/x/kscope/internal/lex"
)

// ErrParse means there is an error in parsing.
var ErrParse = errors.New("parse error")

// errNotDeclaration means the root element is not a declaration.
var errNotDeclaration = errors.New("invalid root token") // NOEXPORT

// parseError holds a parse error and position of the error in the code.
type parseError struct {
	err error
	pos lex.Position
}

// newParseError creates a parseError at provided position p.
func newParseError(p lex.Position, err error) *parseError {
	return &parseError{pos: p, err: err}
}

// Is implements [errors.Is].
func (e *parseError) Is(other error) bool {
	return other == ErrParse
}

// Unwrap implements [errors.Unwrap].
func (e *parseError) Unwrap() error {
	return e.err
}

// Error implements [builtin.error].
func (e *parseError) Error() string {
	return e.err.Error()
}

// Pos returns error's position.
func (e *parseError) Pos() lex.Position {
	return e.pos
}
