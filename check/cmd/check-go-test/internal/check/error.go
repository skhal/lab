// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

// TestError is a test error with the message extracted from the ActionOutput
// events.
type TestError []*TestEvent

// Error implements [builtin.error].
func (err TestError) Error() string {
	buf := new(bytes.Buffer)
	for _, item := range err {
		if item.Action != test.ActionOutput {
			continue
		}
		buf.WriteString(item.Output)
	}
	return strings.TrimRightFunc(buf.String(), unicode.IsSpace)
}

// BuildError is a build error with the message extracted from the ActionOutput
// events.
type BuildError []*BuildEvent

// Error implements [builtin.error].
func (err BuildError) Error() string {
	buf := new(bytes.Buffer)
	for _, item := range err {
		if item.Action != build.ActionOutput {
			continue
		}
		buf.WriteString(item.Output)
	}
	return strings.TrimRightFunc(buf.String(), unicode.IsSpace)
}

// CoverageError means the package has test coverage less than acceptable
// threshold.
type CoverageError struct {
	Package string   // failed package name
	Got     Coverage // actual coverage
	Want    Coverage // threshold
}

// Error implements [builtin.error].
func (err *CoverageError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "=== COVERAGE: %s\n", err.Package)
	fmt.Fprintf(&b, "    coverage: %s of statements\n", err.Got)
	fmt.Fprintf(&b, "    threshold: %s\n", err.Want)
	fmt.Fprintf(&b, "--- FAIL")
	return b.String()
}
