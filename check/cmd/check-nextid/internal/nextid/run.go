// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nextid

// Run checks every file and returns an error on the first failed file.
func Run(files []string) error {
	for _, f := range files {
		if err := CheckFile(f); err != nil {
			return err
		}
	}
	return nil
}
