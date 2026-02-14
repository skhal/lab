// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"errors"
	"os"
)

// Run checks every file and returns a consolidated error of all found
// violations.
func Run(files []string) error {
	var ee []error
	for _, f := range files {
		if err := runFile(f); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func runFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return CheckFile(file, f)
}
