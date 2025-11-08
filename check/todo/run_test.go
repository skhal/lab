// Copyright 2025 Samvel Khalatyan. All rights reserved.

package todo_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/todo"
)

func ExampleRun() {
	readFileFn := func(f string) ([]byte, error) {
		if f != "foo.txt" {
			return nil, fmt.Errorf("error opening file %s", f)
		}
		data := `
// TODO(github.com/foo/bar/issues/123): valid item
// TODO(): invalid item - missing github issue
`
		return []byte(data), nil
	}
	cfg := &todo.Config{
		ReadFileFn: readFileFn,
	}
	if err := todo.Run(cfg, "foo.txt"); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// foo.txt:3 // TODO(): invalid item - missing github issue
	// lint error
}

var TestErr = errors.New("test error")

func TestLinter(t *testing.T) {
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
			name: "valid data violations",
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
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []*todo.Violation
			linter := todo.NewLinter(tc.readFileFn)

			err := linter.Lint(tc.file)
			linter.Visit(func(v *todo.Violation) {
				got = append(got, v)
			})

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("(*todo.Linter).Lint(%q) = %v; want %v", tc.file, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantViolations, got); diff != "" {
				t.Errorf("(*todo.Linter).Visit(%q) mismatch violations (-want, +got):\n%s", tc.file, diff)
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
