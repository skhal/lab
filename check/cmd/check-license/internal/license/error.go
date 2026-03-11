// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound indicates missing license.
	ErrNotFound = errors.New("missing license")

	// ErrInvalid indicates invalid license.
	ErrInvalid = errors.New("invalid license")

	// ErrBinaryFile indicates the file is binary
	ErrBinaryFile = errors.New("binary file")
)

// InvalidError is [ErrInvalid] with line number where copyright was found.
type InvalidError struct {
	line int
}

// NewInvalidError returns a new InvalidError for provided line number.
func NewInvalidError(line int) *InvalidError {
	return &InvalidError{
		line: line,
	}
}

// Error implements [builtin.error] interface.
func (e *InvalidError) Error() string {
	return fmt.Sprintf("L%d: %s", e.line, ErrInvalid)
}

// Is equates InvalidError to ErrInvalid.
func (e *InvalidError) Is(target error) bool {
	return target == ErrInvalid
}
