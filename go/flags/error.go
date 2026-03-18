// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags

import (
	"errors"
	"flag"
	"fmt"
)

// ErrFlag means the flag has invalid value.
var ErrFlag = errors.New("invalid flag")

// FlagError is an error attached to the flag.
type FlagError struct {
	f   *flag.Flag
	err error
}

// NewFlagError creates a [FlagError].
func NewFlagError(f *flag.Flag, err error) error {
	return &FlagError{f, err}
}

// Error implements [builtin.Error] interface.
func (e FlagError) Error() string {
	return fmt.Sprintf("%s %s: %v, %s", ErrFlag, e.f.Name, e.err, e.f.Value)
}

// Is treats FlagError equivalent to [ErrFlag].
func (e FlagError) Is(err error) bool {
	return err == ErrFlag
}

// Unwrap returns the error that is associated with the flag.
func (e FlagError) Unwrap() error {
	return e.err
}
