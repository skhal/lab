// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"errors"
	"fmt"
	"os"
)

var (
	// ErrNotFound indicates missing license.
	ErrNotFound = errors.New("missing license")

	// ErrInvalid indicates invalid license.
	ErrInvalid = errors.New("invalid license")
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

func (e *InvalidError) Error() string {
	return fmt.Sprintf("L%d: %s", e.line, ErrInvalid)
}

func (e *InvalidError) Is(target error) bool {
	return target == ErrInvalid
}

// Run checks that every file in files has a license.
func Run(files []string) error {
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if err := Check(data); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
	}
	return nil
}
