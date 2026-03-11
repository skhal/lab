// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

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
	return doRun(files, opts)
}

var runTimeout = 100 * time.Millisecond

func doRun(files []string, opts *RunOptions) error {
	errch := make(chan error)
	go func() {
		defer close(errch)
		ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
		defer cancel()
		var wg sync.WaitGroup
		for _, f := range files {
			wg.Add(1)
			go func(ctx context.Context, f string) {
				defer wg.Done()
				err := check(f, opts)
				if errors.Is(err, ErrBinaryFile) {
					fmt.Fprintln(os.Stderr, err, ": skip")
				} else if err != nil {
					select {
					case errch <- err:
					case <-ctx.Done():
					}
				}
			}(ctx, f)
		}
		wg.Wait()
	}()
	var ee []error
	for err := range errch {
		ee = append(ee, err)
	}
	return errors.Join(ee...)
}

func check(file string, opts *RunOptions) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if isBinary(data) {
		return fmt.Errorf("check %s: %w", file, ErrBinaryFile)
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

const nullByte = 0

func isBinary(b []byte) bool {
	return bytes.IndexByte(b, nullByte) != -1
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
