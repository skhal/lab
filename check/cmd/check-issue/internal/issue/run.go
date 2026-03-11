// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package issue

import (
	"bytes"
	"errors"
	"os"
	"regexp"
)

// ErrNoIssue means missing issue reference.
var ErrNoIssue = errors.New("missing issue")

// Run verifies that the files has an issue reference.
func Run(files []string) error {
	var ee []error
	for _, f := range files {
		if err := run(f); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func run(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return Check(b)
}

var (
	reNoIssue = regexp.MustCompile(`^(?i)no_issue(?:: .*)?$`)
	reIssue   = regexp.MustCompile(`^(?i)(?:close|fix|issue) (?:\w+/\w+)?#\d+$`)
)

const eol = "\n"

// Check verifies that the buffer b has an valid issue reference. A valid issue
// reference (ignore case) is separate on the line and has one of the forms:
//
//	no_issue
//	verb reference
//
// where verb is one of `issue`, `close`, `fix` and reference is either `#123`
// or `owner/repo#123`.
//
// Examples:
//
//	NO_ISSUE
//	Issue #123
//	Close user/repo#123
func Check(b []byte) error {
	for ln := range bytes.Lines(b) {
		ln = bytes.TrimRight(ln, eol)
		if validate(ln) {
			return nil
		}
	}
	return ErrNoIssue
}

func validate(b []byte) bool {
	if reNoIssue.Match(b) {
		return true
	}
	if reIssue.Match(b) {
		return true
	}
	return false
}
