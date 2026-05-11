// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"path/filepath"
	"slices"
	"strings"

	goslices "github.com/skhal/lab/go/slices"
)

// Run runs `go test` on packages for listed files.
func Run(files []string, opts ...Opt) error {
	packages := CollectPackages(files)
	if len(packages) == 0 {
		return nil
	}
	return NewTester(opts...).Test(packages)
}

// CollectPackages collects a unique set of packages from a list of files with
// .go extension.
func CollectPackages(files []string) []string {
	packages := slices.Collect(Unique(Paths(FilterFunc(slices.Values(files), IsGoFile))))
	packages = slices.DeleteFunc(packages, func(p string) bool {
		return strings.Contains(p, "/testdata/")
	})
	packages = goslices.MapFunc(packages, func(p string) string {
		// Gotest is a pre-commit check. Git reports changed files with respect
		// to the work tree, without leading "./". `go test` expects local packages
		// to start with "./".
		return filepath.FromSlash("./" + filepath.Clean(p))
	})
	return packages
}

// IsGoFile reports whether a file is a Go source, i.e. has `.go` extension.
func IsGoFile(f string) bool {
	const goExtension = ".go"
	return filepath.Ext(f) == goExtension
}
