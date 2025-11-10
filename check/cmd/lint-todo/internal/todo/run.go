// Copyright 2025 Samvel Khalatyan. All rights reserved.

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

// ErrLint indicates an error in linting inputs.
var ErrLint = errors.New("lint error")

var (
	todoPrefix = regexp.MustCompile(`(?i)(?://|#|")\s+\btodo\b.*`)
	todoRegexp = regexp.MustCompile(`TODO\(github.com/\w+/\w+/issues/\d+\):\s.+$`)
)

var (
	nolintRegexp = regexp.MustCompile(`^(?://|#|") lint-todo off(?:: .*)?$`)
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
	l := NewLinter(cfg.ReadFileFn)
	for _, file := range files {
		if err := l.Lint(file); err != nil {
			return err
		}
	}
	if !l.HasViolations() {
		return nil
	}
	l.Visit(func(v *Violation) { fmt.Println(v) })
	return ErrLint
}

// Linter implements todo-line validation logic.
type Linter struct {
	readFileFn ReadFileFunc
	vv         []*Violation
}

// NewLinter creates a Linter.
func NewLinter(f ReadFileFunc) *Linter {
	return &Linter{
		readFileFn: f,
	}
}

// HasViolations reports whether violations were found.
func (l *Linter) HasViolations() bool {
	return len(l.vv) != 0
}

// Lint validates file. It reutns an error if reading the file fails. Linter
// keeps track of found violations. Use (*Linter).HasViolations() to check
// for findings and (*linter).Visit() to access violations.
func (l *Linter) Lint(file string) error {
	data, err := l.readFileFn(file)
	if err != nil {
		return err
	}
	for v := range findViolations(data) {
		v.File = file
		l.vv = append(l.vv, v)
	}
	return nil
}

// Visit calls fn for every found violation.
func (l *Linter) Visit(fn func(v *Violation)) {
	for _, v := range l.vv {
		fn(v)
	}
}

func findViolations(data []byte) iter.Seq[*Violation] {
	return func(yield func(*Violation) bool) {
		s := bufio.NewScanner(bytes.NewBuffer(data))
		for row := 1; s.Scan(); row += 1 {
			line := s.Text()
			if nolintRegexp.MatchString(line) {
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
