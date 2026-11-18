// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-yaml-lint runs Ansible linter on Ansible files and yamlfmt otherwise.
//
// # Synopsis
//
//	check-yaml-lint file [file...]
//
// # Description
//
// check-yaml-lint groups files by Ansible project, collection, and everything
// else. It runs ansible-lint for the first two kinds and yamlfmt for the
// latter.
//
// The check detects Ansible project by looking for ansible.cfg in one of the
// parent folders. Ansible collection is a folder with galaxy.yml file.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/check/cmd/check-yaml-lint/internal/check"
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s file ...\n", filepath.Base(os.Args[0]))
	}
}

func main() {
	flag.Parse()
	if err := check.Run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
