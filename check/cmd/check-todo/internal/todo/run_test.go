// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// check-todo off

package todo_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

var TestErr = errors.New("test error")

func TestChecker(t *testing.T) {
	tests := []struct {
		name           string
		readFileFn     todo.ReadFileFunc
		file           string
		wantErr        error
		wantViolations []*todo.Violation
	}{
		{
			name:       "error reading file",
			readFileFn: func(string) ([]byte, error) { return nil, TestErr },
			file:       "test.txt",
			wantErr:    TestErr,
		},
		{
			name: "empty data",
			readFileFn: func(string) ([]byte, error) {
				return []byte(``), nil
			},
			file: "test.txt",
		},
		{
			name: "no violations",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`// TODO(github.com/foo/bar/issues/123): test`), nil
			},
			file: "test.txt",
		},
		{
			name: "missing issue",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`// TODO(): test`), nil
			},
			file: "test.txt",
			wantViolations: []*todo.Violation{
				makeViolation("test.txt", 1, `// TODO(): test`),
			},
		},
		{
			name: "missing description",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`// TODO(github.com/foo/bar/issues/123)`), nil
			},
			file: "test.txt",
			wantViolations: []*todo.Violation{
				makeViolation("test.txt", 1, `// TODO(github.com/foo/bar/issues/123)`),
			},
		},
		{
			name: "multiple violations",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`// TODO(github.com/foo/bar/issues/123)
// TODO(): test
				`), nil
			},
			file: "test.txt",
			wantViolations: []*todo.Violation{
				makeViolation("test.txt", 1, `// TODO(github.com/foo/bar/issues/123)`),
				makeViolation("test.txt", 2, `// TODO(): test`),
			},
		},
		{
			name: "disable lint on multiple violations",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`// check-todo off
// TODO(github.com/foo/bar/issues/123)
// TODO(): test
				`), nil
			},
			file: "test.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []*todo.Violation
			checker := todo.NewChecker(tc.readFileFn)

			err := checker.Check(tc.file)
			checker.Visit(func(v *todo.Violation) {
				got = append(got, v)
			})

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("(*todo.Checker).Check(%q) = %v; want %v", tc.file, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantViolations, got); diff != "" {
				t.Errorf("(*todo.Checker).Visit(%q) mismatch violations (-want, +got):\n%s", tc.file, diff)
			}
		})
	}
}

func makeViolation(f string, row int, s string) *todo.Violation {
	return &todo.Violation{
		File: f,
		Row:  row,
		Line: s,
	}
}
