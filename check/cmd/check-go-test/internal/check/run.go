// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
)

// ErrTest indicates an error in running go tests.
var ErrTest = errors.New("test error")

// Run runs `go test` on packages for listed files.
func Run(files ...string) (err error) {
	packages := slices.Collect(Unique(Paths(FilterFunc(slices.Values(files), IsGoFile))))
	if len(packages) == 0 {
		return
	}
	tester := NewTester()
	for _, p := range packages {
		// Gotest is a pre-commit check. Git reports changed files with respect
		// to the work tree, without leading "./". `go test` expects local packages
		// to start with "./".
		p = filepath.FromSlash("./" + filepath.Clean(p))
		if err = tester.Test(p); err != nil {
			return
		}
	}
	tester.VisitFails(func(f *FailedTest) {
		err = ErrTest
		fmt.Print(string(f.Output))
	})
	return
}

// IsGoFile reports whether a file is a Go source, i.e. has `.go` extension.
func IsGoFile(f string) bool {
	const goExtension = ".go"
	return filepath.Ext(f) == goExtension
}
