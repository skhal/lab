// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import "os"

// Files is a collection of file paths.
type Files []string

// Run groups files by Ansible project, collection, and everything. It passes
// the files through ansible-lint for the first two groups and yamlfmt for the
// last one.
func Run(files []string) error {
	fs := os.DirFS(".")
	ansibleProjects, nonAnsibleFiles := Collect(fs, files)
	return Check(ansibleProjects, nonAnsibleFiles)
}
