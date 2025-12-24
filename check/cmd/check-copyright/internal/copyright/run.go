// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package copyright

import (
	"errors"
	"regexp"
)

// ErrNotFound indicates missing or invalid copyright.
var ErrNotFound = errors.New("copyright is not found")

var (
	reCopyright = regexp.MustCompile(`(//|#|") Copyright \d{4} \w+( \w+)?. All rights reserved.
(//|#|")
(//|#|") Use of this source code is governed by a BSD-style
(//|#|") license that can be found in the LICENSE file.`)
)

// ReadFileFn reads file and returns its content or error.
type ReadFileFn = func(string) ([]byte, error)

// Config configures the runner.
type Config struct {
	ReadFile ReadFileFn
}

// Run checks whether the file include a copyright comment.
func Run(cfg *Config, files []string) error {
	for _, f := range files {
		data, err := cfg.ReadFile(f)
		if err != nil {
			return err
		}
		if err := check(data); err != nil {
			return err
		}
	}
	return nil
}

func check(buf []byte) error {
	match := reCopyright.Find(buf)
	if match == nil {
		return ErrNotFound
	}
	return nil
}
