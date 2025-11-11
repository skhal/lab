// Copyright 2025 Samvel Khalatyan. All rights reserved.

// Package nosubmit lints files against presence of "DO NOT SUBMIT" comment.
//
// Place "DO NOT SUBMIT" comment anywhere in the code to prevent it from being
// committed. Add optional descriptin.
//
// Example:
//   # DO NOT SUBMIT: work in progress
package nosubmit

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"iter"
	"regexp"
)

// ErrLint indicates presense of "DO NO SUBMIT" comment.
var ErrLint = errors.New("lint error")

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
		if Lint(data) {
			return ErrLint
		}
	}
	return nil
}

var nosubmitRe = regexp.MustCompile(`(?i)(?://|#|")\s+DO NOT SUBMIT(?:: .*)?$`)

// Lint checks if bytes b include "DO NOT SUBMIT" comment.
func Lint(b []byte) bool {
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
