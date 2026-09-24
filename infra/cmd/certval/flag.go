// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
)

// FlagError helps distinguish [flag] package errors.
type FlagError struct {
	err error
}

// Error implements [builtin.error] interface.
func (e FlagError) Error() string {
	return e.err.Error()
}

// Unwrap returns wrapped error for [errors.Unwrap].
func (e FlagError) Unwrap() error {
	return e.err
}

type flagValidatorFunc func() error

// ParseFlags parses flags.
func ParseFlags(validator flagValidatorFunc) error {
	flag.Parse()
	if err := validator(); err != nil {
		return FlagError{err}
	}
	return nil
}
