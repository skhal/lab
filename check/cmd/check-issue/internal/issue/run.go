// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package issue

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
)

var (
	// ErrCheck indicates general error in the check.
	ErrCheck = errors.New("check error")

	// ErrNoIssue indicates missing issue
	ErrNoIssue = errors.New("missing issue")
)

// Run executes the check. It is expected that the check will run as a
// commit-msg git-hook(1). Therefore there should be a single file, else it
// returns an error.
func Run(files []string) error {
	if len(files) != 1 {
		return ErrCheck
	}
	return run(files[0])
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
	reIssue   = regexp.MustCompile(`^(?i)(?:issue|close|fix) (?:\w+/\w+)?#\d+$`)
)

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
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		line := s.Text()
		if reNoIssue.MatchString(line) {
			return nil
		}
		if reIssue.MatchString(line) {
			return nil
		}
	}
	return ErrNoIssue
}
