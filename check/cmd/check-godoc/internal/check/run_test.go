// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-godoc/internal/check"
)

func TestPathCollector(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		wantPaths []string
	}{
		{
			name:      "one file",
			files:     []string{"/a/1.go"},
			wantPaths: []string{"/a"},
		},
		{
			name:      "two files",
			files:     []string{"/a/1.go", "/b/2.go"},
			wantPaths: []string{"/a", "/b"},
		},
		{
			name:      "two files same dir",
			files:     []string{"/a/1.go", "/a/2.go"},
			wantPaths: []string{"/a"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pc := check.NewPathCollector()

			for _, f := range tc.files {
				pc.CollectFile(f)
			}

			if diff := cmp.Diff(tc.wantPaths, pc.Paths()); diff != "" {
				t.Errorf("(check.PathCollector).Paths() mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestIsTest(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "foo.go"},
		{name: "foo_test.go", want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := check.IsTest(tc.name)

			if got != tc.want {
				t.Errorf("check.IsTest(%q) got %v; want %v", tc.name, got, tc.want)
			}
		})
	}
}
