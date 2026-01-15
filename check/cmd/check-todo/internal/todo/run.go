// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package todo validates todo-lines.
package todo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"iter"
	"os"
	"regexp"
)

// ErrCheck indicates an error in the input.
var ErrCheck = errors.New("check error")

var (
	todoPrefix = regexp.MustCompile(`(?i)(?://|#|")\s+\btodo\b.*`)
	todoRegexp = regexp.MustCompile(`TODO\(github.com/\w+/\w+/issues/\d+\):\s.+$`)
)

var (
	nocheckRegexp = regexp.MustCompile(`^(?://|#|") check-todo off(?:: .*)?$`)
)

// ReadFileFunc reads file and returns data or error.
type ReadFileFunc func(string) ([]byte, error)

// Config holds parameters for todo-linter.
type Config struct {
	// ReadFileFn reads file contents.
	ReadFileFn ReadFileFunc
}

// NewConfig constructructs a default configuration with os.ReadFile to read
// files.
func NewConfig() *Config {
	return &Config{
		ReadFileFn: os.ReadFile,
	}
}

// Run validates todo-lines in files.
func Run(cfg *Config, files ...string) (err error) {
	chk := NewChecker(cfg.ReadFileFn)
	for _, file := range files {
		if err := chk.Check(file); err != nil {
			return err
		}
	}
	if !chk.HasViolations() {
		return nil
	}
	chk.Visit(func(v *Violation) { fmt.Println(v) })
	return ErrCheck
}

// Checker implements todo-line validation logic.
type Checker struct {
	readFileFn ReadFileFunc
	vv         []*Violation
}

// NewChecker creates a Checker.
func NewChecker(f ReadFileFunc) *Checker {
	return &Checker{
		readFileFn: f,
	}
}

// HasViolations reports whether violations were found.
func (chk *Checker) HasViolations() bool {
	return len(chk.vv) != 0
}

// Check validates file. It reutns an error if reading the file fails. Checker
// keeps track of found violations. Use (*Checker).HasViolations() to check
// for findings and (*linter).Visit() to access violations.
func (chk *Checker) Check(file string) error {
	data, err := chk.readFileFn(file)
	if err != nil {
		return err
	}
	for v := range findViolations(data) {
		v.File = file
		chk.vv = append(chk.vv, v)
	}
	return nil
}

// Visit calls fn for every found violation.
func (chk *Checker) Visit(fn func(v *Violation)) {
	for _, v := range chk.vv {
		fn(v)
	}
}

func findViolations(data []byte) iter.Seq[*Violation] {
	return func(yield func(*Violation) bool) {
		s := bufio.NewScanner(bytes.NewBuffer(data))
		for row := 1; s.Scan(); row += 1 {
			line := s.Text()
			if nocheckRegexp.MatchString(line) {
				break
			}
			if !todoPrefix.MatchString(line) {
				continue
			}
			if todoRegexp.MatchString(line) {
				continue
			}
			v := &Violation{
				Row:  row,
				Line: line,
			}
			if !yield(v) {
				break
			}
		}
	}
}

// Violation holds invalid line and its location in the file.
type Violation struct {
	File string
	Row  int
	Line string
}

// String implements fmt.Stringer interface.
func (v *Violation) String() string {
	return fmt.Sprintf("%s:%d %s", v.File, v.Row, v.Line)
}
