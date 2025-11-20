// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package issue_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/skhal/lab/check/cmd/check-issue/internal/issue"
)

func ExampleRun() {
	readFileFn := func(f string) ([]byte, error) {
		data := `
Issue #123
`
		return []byte(data), nil
	}
	cfg := &issue.Config{
		ReadFileFn: readFileFn,
	}
	if err := issue.Run(cfg, []string{"foo.txt"}); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
}

var TestErr = errors.New("test error")

func TestCheck(t *testing.T) {
	tests := []struct {
		name       string
		readFileFn issue.ReadFileFunc
		file       string
		wantErr    error
	}{
		{
			name:       "error readinf eil",
			readFileFn: func(string) ([]byte, error) { return nil, TestErr },
			file:       "test.txt",
			wantErr:    TestErr,
		},
		{
			name: "missing issue",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Test commit`), nil
			},
			file:    "test.txt",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "no issue tag",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`NO_ISSUE`), nil
			},
			file: "test.txt",
		},
		{
			name: "no issue tag with description",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`NO_ISSUE: N/A`), nil
			},
			file: "test.txt",
		},
		{
			name: "issue local",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Issue #123`), nil
			},
			file: "test.txt",
		},
		{
			name: "issue other owher",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Issue owner/repo#123`), nil
			},
			file: "test.txt",
		},
		{
			name: "close local",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Close #123`), nil
			},
			file: "test.txt",
		},
		{
			name: "close other owher",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Close owner/repo#123`), nil
			},
			file: "test.txt",
		},
		{
			name: "fix local",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Fix #123`), nil
			},
			file: "test.txt",
		},
		{
			name: "fix other owher",
			readFileFn: func(string) ([]byte, error) {
				return []byte(`Fix owner/repo#123`), nil
			},
			file: "test.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &issue.Config{
				ReadFileFn: tc.readFileFn,
			}

			err := issue.Check(cfg, tc.file)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("issue.Check(_, %q) = %v; want error %v", tc.file, err, tc.wantErr)
			}
		})
	}
}
