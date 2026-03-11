// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package todo validates todo-lines.
package todo

import (
	"bytes"
	"errors"
	"os"
	"regexp"
)

var (
	reNoLint = regexp.MustCompile(`^(?://|#|") check-todo off(?:: .*)?$`)
	rePrefix = regexp.MustCompile(`(?i)(?://|#|")\s+\btodo\b.*`)
	reTodo   = regexp.MustCompile(`TODO\(github.com/\w+/\w+/issues/\d+\):\s.+$`)
)

// Run checks todo-comments in files. It collects found errors from files and
// reports an error with all found violations if any.
func Run(files []string) (err error) {
	var ee []error
	for _, f := range files {
		if err := run(f); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func run(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return NewChecker(file).Check(b)
}

// Checker checks todo-comments in the file.
type Checker struct {
	file string
}

// NewChecker creates a file checker.
func NewChecker(file string) *Checker {
	return &Checker{file}
}

const eol = "\n"

// Check validates todo-comments in b.
func (ch *Checker) Check(b []byte) error {
	var ee []error
	var line int
	for lb := range bytes.Lines(b) {
		line++
		lb = bytes.TrimRight(lb, eol)
		if reNoLint.Match(lb) {
			break
		}
		if !rePrefix.Match(lb) {
			continue
		}
		if reTodo.Match(lb) {
			continue
		}
		ee = append(ee, &TodoError{
			File: ch.file,
			Line: line,
			Text: string(lb),
		})
	}
	return errors.Join(ee...)
}
