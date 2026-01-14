// Copyright 2025 Samvel Khalatyan. All rights reserved.
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

// ReadFileFn reads file and returns its content or error.
type ReadFileFn = func(string) ([]byte, error)

// Config configures the runner.
type Config struct {
	ReadFile ReadFileFn
}

// Run checks whether the file include a copyright comment.
func Run(cfg *Config, files []string) error {
	for _, f := range files {
		data, err := cfg.ReadFile(f)
		if err != nil {
			return err
		}
		if err := Check(data); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
	}
	return nil
}
