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

// Run checks that every file in files has a license.
func Run(files []string, fix bool, holder string) error {
	errch := make(chan error)
	go func() {
		defer close(errch)
		ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
		defer cancel()
		var wg sync.WaitGroup
		for _, f := range files {
			wg.Go(func() {
				switch err := check(f, fix, holder); {
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

func check(file string, fix bool, holder string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if isBinary(data) {
		return fmt.Errorf("check %s: %w", file, ErrBinaryFile)
	}
	err = Check(data)
	if err != nil {
		if fix && errors.Is(err, ErrNotFound) {
			return runFix(file, data, holder)
		}
		return fmt.Errorf("%s: %w", file, err)
	}
	return nil
}

const nullByte = 0

func isBinary(b []byte) bool {
	return bytes.IndexByte(b, nullByte) != -1
}

func runFix(file string, data []byte, holder string) error {
	data, err := Add(data, file, holder)
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
