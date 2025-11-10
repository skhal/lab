// Copyright 2025 Samvel Khalatyan. All rights reserved.

package test

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
	r := NewTester()
	for _, p := range packages {
		// Gotest is a pre-commit check. Git reports changed files with respect
		// to the wortree, without leading "./". `go test` assumes that packages
		// without "./" prefix are non-local.
		p = filepath.FromSlash("./" + p)
		if err = r.Test(p); err != nil {
			return
		}
	}
	r.VisitFails(func(f *FailedTest) {
		err = ErrTest
		fmt.Print(string(f.Output))
	})
	return
}

// IsGoFile reports whether a file has .go extension.
func IsGoFile(f string) bool {
	const goExtension = ".go"
	return filepath.Ext(f) == goExtension
}
