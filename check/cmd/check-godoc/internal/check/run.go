// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check package ensures that Go exported symbols and the package include
// documentation.
package check

import (
	"errors"
	"path/filepath"
	"strings"
)

// Run verifies that every non-generated Go file has documentation attached to
// the exported declarations. It returns an error on the first failed file.
func Run(files []string) error {
	var ee []error
	pc := NewPathCollector()
	for _, f := range files {
		if IsTest(f) {
			continue
		}
		if err := CheckFile(f); err != nil {
			ee = append(ee, err)
		}
		pc.CollectFile(f)
	}
	for _, d := range pc.Paths() {
		if err := CheckDir(d); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

// PathCollector collects unique folders from a set of files.
type PathCollector struct {
	seen  map[string]bool
	paths []string
}

// NewPathCollector creates a new PathCollector.
func NewPathCollector() *PathCollector {
	return &PathCollector{
		seen: make(map[string]bool),
	}
}

// CollectFile collects the path to the file. It skips the path if already in
// the collector.
func (pc *PathCollector) CollectFile(name string) {
	d := filepath.Dir(name)
	if pc.seen[d] {
		return
	}
	pc.seen[d] = true
	pc.paths = append(pc.paths, d)
}

// Paths returns a slice of collected paths.
func (pc *PathCollector) Paths() []string {
	return pc.paths
}

// IsTest reports whether the file is a test file. A test file has _test.go
// suffix.
func IsTest(name string) bool {
	return strings.HasSuffix(name, "_test.go")
}
