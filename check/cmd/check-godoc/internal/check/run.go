// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check doc
package check

import (
	"strings"
)

// Run verifies that every non-generated Go file has documentation attached to
// the exported declarations. It returns an error on the first failed file.
func Run(files []string) error {
	for _, f := range files {
		if IsTest(f) {
			continue
		}
		if err := Check(f); err != nil {
			return err
		}
	}
	return nil
}

// IsTest reports whether the file is a test file. A test file has _test.go
// suffix.
func IsTest(name string) bool {
	return strings.HasSuffix(name, "_test.go")
}
