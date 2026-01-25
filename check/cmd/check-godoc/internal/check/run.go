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
	seen := make(map[string]bool)
	dirs := make([]string, 0, len(files))
	for _, f := range files {
		if IsTest(f) {
			continue
		}
		if err := CheckFile(f); err != nil {
			ee = append(ee, err)
		}
		d := filepath.Dir(f)
		if seen[d] {
			continue
		}
		seen[d] = true
		dirs = append(dirs, d)
	}
	for _, d := range dirs {
		if err := CheckDir(d); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

// IsTest reports whether the file is a test file. A test file has _test.go
// suffix.
func IsTest(name string) bool {
	return strings.HasSuffix(name, "_test.go")
}
