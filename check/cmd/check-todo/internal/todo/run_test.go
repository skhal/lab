// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// check-todo off

package todo_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

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
			name: "no issue",
			file: "test",
			s:    "// TODO(): test",
			wantErr: &todo.TodoError{
				File: "test",
				Line: 1,
				Text: "// TODO(): test",
			},
		},
		{
			name: "no description",
			file: "test",
			s:    "// TODO(github.com/foo/bar/issues/123)",
			wantErr: &todo.TodoError{
				File: "test",
				Line: 1,
				Text: "// TODO(github.com/foo/bar/issues/123)",
			},
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
