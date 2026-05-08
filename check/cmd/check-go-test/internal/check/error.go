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

type testError []Event

// Error implements [builtin.error].
func (err testError) Error() string {
	buf := new(bytes.Buffer)
	for _, item := range err {
		te := item.(*TestEvent)
		if te.Action != test.ActionOutput {
			continue
		}
		buf.WriteString(te.Output)
	}
	return strings.TrimRightFunc(buf.String(), unicode.IsSpace)
}

type buildError []Event

// Error implements [builtin.error].
func (err buildError) Error() string {
	buf := new(bytes.Buffer)
	for _, item := range err {
		be := item.(*BuildEvent)
		if be.Action != build.ActionOutput {
			continue
		}
		buf.WriteString(be.Output)
	}
	return strings.TrimRightFunc(buf.String(), unicode.IsSpace)
}

type coverageError struct {
	pkg  string
	got  Coverage
	want Coverage
}

// Error implements [builtin.error].
func (err *coverageError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "=== COVERAGE: %s\n", err.pkg)
	fmt.Fprintf(&b, "    coverage: %s of statements\n", err.got)
	fmt.Fprintf(&b, "    threshold: %s\n", err.want)
	fmt.Fprintf(&b, "--- FAIL")
	return b.String()
}
