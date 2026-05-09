// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// check-todo off

package todo_test

import (
	"errors"
	"os"
	"testing"

	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr error
	}{
		{
			name:    "file not exist",
			file:    "testdata/not_exist_file.txt",
			wantErr: os.ErrNotExist,
		},
		{
			name: "file without todo",
			file: "testdata/no_todo.txt",
		},
		{
			name:    "file with todo",
			file:    "testdata/todo.txt",
			wantErr: todo.ErrTodo,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := todo.Run([]string{tc.file})

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error %v; want %v", err, tc.wantErr)
			}
		})
	}
}

func TestChecker(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		s       string
		wantErr error
	}{
		{
			name: "empty data",
			file: "test",
		},
		{
			name: "valid",
			file: "test",
			s:    "// TODO(github.com/foo/bar/issues/123): test",
		},
		{
			name:    "no issue",
			file:    "test",
			s:       "// TODO(): test",
			wantErr: todo.ErrTodo,
		},
		{
			name:    "no description",
			file:    "test",
			s:       "// TODO(github.com/foo/bar/issues/123)",
			wantErr: todo.ErrTodo,
		},
		{
			name: "no lint",
			file: "test",
			s: `
// check-todo off
no issue
// TODO(): test
no description
// TODO(github.com/foo/bar/issues/123)
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ch := todo.NewChecker(tc.file)

			err := ch.Check([]byte(tc.s))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Check() error mismatch.\ngot: %v\nwant: %v", err, tc.wantErr)
				t.Logf("\nfile: %q\ndata:\n%s", tc.file, tc.s)
			}
		})
	}
}
