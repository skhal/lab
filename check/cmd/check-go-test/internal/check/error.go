// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

var (
	// ErrBuild means there is an error in building a test package.
	ErrBuild = errors.New("build error")

	// ErrCoverage means the test coverage is below set threshold.
	ErrCoverage = errors.New("coverage error")

	// ErrTest means there is an error in one of the test cases.
	ErrTest = errors.New("test error")
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

// Is makes TestError equivalent to ErrTest.
func (TestError) Is(err error) bool {
	return err == ErrTest
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

// Is makes BuildError equivalent to ErrBuild.
func (BuildError) Is(err error) bool {
	return err == ErrBuild
}

// CoverageError means the package has test coverage less than acceptable
// threshold.
type CoverageError struct {
	Package string   // failed package name
	Got     Coverage // actual coverage
	Want    Coverage // threshold
}

// Is makes CoverageError equivalent to ErrCoverage.
func (*CoverageError) Is(err error) bool {
	return err == ErrCoverage
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
