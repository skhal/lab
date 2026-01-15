// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	// ErrNotFound indicates missing license.
	ErrNotFound = errors.New("missing license")

	// ErrInvalid indicates invalid license.
	ErrInvalid = errors.New("invalid license")

	// ErrFlags indicates error with input flags.
	ErrFlags = errors.New("invalid flags")
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

// RunOptions customize the license check.
type RunOptions struct {
	Fix    bool   // fix license if missing
	Holder string // attribute license to the holder
}

// RegisterFlags registers RunOptions with flags.
func (opt *RunOptions) RegisterFlags(fs *flag.FlagSet) {
	fs.BoolVar(&opt.Fix, "fix", false, "insert license if missing")
	fs.StringVar(&opt.Holder, "holder", "", "license holder name")
}

// Validate ensures that flag dependencies are satisfied. For example,
// [RunOptions.Holder] must be set along with [RunOPtions.Fix].
func (opt *RunOptions) Validate() error {
	if opt.Fix && opt.Holder == "" {
		return fmt.Errorf("%w: missing -holder with -fix", ErrFlags)
	}
	return nil
}

// Run checks that every file in files has a license.
func Run(files []string, opts *RunOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	for _, f := range files {
		err := check(f, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func check(file string, opts *RunOptions) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = Check(data)
	if err != nil {
		if opts.Fix && errors.Is(err, ErrNotFound) {
			return fix(file, data, opts)
		}
		return fmt.Errorf("%s: %w", file, err)
	}
	return nil
}

func fix(file string, data []byte, opts *RunOptions) error {
	data, err := Add(data, file, opts.Holder)
	if err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, data, info.Mode())
	if err != nil {
		return err
	}
	// use error to indicate the file change
	return fmt.Errorf("%s: fixed file", file)
}
