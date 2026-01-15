// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Check-issue verifies that git commit message includes an issue.

Use it as a commit-msg git-hook(1). Supported formats are:
- `Issue #123`
- `NO_ISSUE`
- `NO_ISSUE: description`
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-issue/internal/issue"
)

func main() {
	flag.Parse()
	cfg := issue.NewConfig()
	if err := issue.Run(cfg, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
