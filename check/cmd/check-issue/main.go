// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-issue verifies that a git-commit(1) message includes a reference to
// an issue or explicitly states that the commit is not associated with any
// issue.
//
// The issue reference is case insensitive and can be in one of the following
// forms:
//
//	NO_ISSUE
//	NO_ISSUE: description
//	Issue #123
//	Fix #123
//	Close #123
//
// The issue number can also include the owher and repo, e.g. `owner/repo#123`.
//
// SYNOPSIS
//
//	check-issue .git/COMMIT_EDITMSG
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-issue/internal/issue"
)

func main() {
	if err := issue.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
