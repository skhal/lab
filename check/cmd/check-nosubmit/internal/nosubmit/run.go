// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package nosubmit checks files against presence of "DO NOT SUBMIT" comment.
//
// Place "DO NOT SUBMIT" comment anywhere in the code to prevent it from being
// committed. Add optional description.
//
// Example:
//
//	# DO NOT SUBMIT: work in progress
package nosubmit

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"iter"
	"regexp"
)

// ErrCheck indicates presence of "DO NO SUBMIT" comment.
var ErrCheck = errors.New("check error")

// ReadFileFunc reads file and returns contents or error.
type ReadFileFunc func(string) ([]byte, error)

// Config passes ReadFileFunc.
type Config struct {
	ReadFileFn ReadFileFunc
}

// Run checks whether any of the files include "DO NOT SUBMIT" comment.
func Run(ctx context.Context, cfg *Config, files ...string) error {
	for _, f := range files {
		data, err := cfg.ReadFileFn(f)
		if err != nil {
			return err
		}
		if Check(data) {
			return ErrCheck
		}
	}
	return nil
}

var nosubmitRe = regexp.MustCompile(`(?i)(?://|#|")\s+DO NOT SUBMIT(?:: .*)?$`)

// Check checks if bytes b include "DO NOT SUBMIT" comment.
func Check(b []byte) bool {
	for line := range scanLines(b) {
		if nosubmitRe.Match(line) {
			return true
		}
	}
	return false
}

func scanLines(b []byte) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		scanner := bufio.NewScanner(bytes.NewReader(b))
		for scanner.Scan() {
			if !yield(scanner.Bytes()) {
				break
			}
		}
	}
}
