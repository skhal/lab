// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var runTimeout = 100 * time.Millisecond

// Run checks copyright block in the files. It returns an error with a list of
// files that miss or have invalid license.
func Run(files []string) error {
	return run(files, new(noopFixer))
}

// Fix checks copyright blocks in the files and fixes the files if the block
// is missing.
func Fix(files []string, holder string) error {
	return run(files, &copyrightFixer{holder: holder})
}

type fixer interface {
	Fix(file string, b []byte) error
}

type noopFixer struct{}

// Fix returns an error about missing copyright.
func (fx *noopFixer) Fix(file string, _ []byte) error {
	return fmt.Errorf("%s: %w", file, ErrNotFound)
}

type copyrightFixer struct {
	holder string
}

// Fix adds a copyright block to the file content.
func (fx *copyrightFixer) Fix(file string, data []byte) error {
	data, err := Add(data, file, fx.holder)
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

func run(files []string, fx fixer) error {
	errch := make(chan error)
	go func() {
		defer close(errch)
		ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
		defer cancel()
		var wg sync.WaitGroup
		for _, f := range files {
			wg.Go(func() {
				switch err := check(f, fx); {
				case errors.Is(err, ErrBinaryFile):
					// skip
				case err != nil:
					select {
					case errch <- err:
					case <-ctx.Done():
					}
				}
			})
		}
		wg.Wait()
	}()
	var ee []error
	for err := range errch {
		ee = append(ee, err)
	}
	return errors.Join(ee...)
}

func check(file string, fx fixer) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if isBinary(data) {
		return fmt.Errorf("check %s: %w", file, ErrBinaryFile)
	}
	err = Check(data)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fx.Fix(file, data)
		}
		return fmt.Errorf("%s: %w", file, err)
	}
	return nil
}

const nullByte = 0

func isBinary(b []byte) bool {
	return bytes.IndexByte(b, nullByte) != -1
}
