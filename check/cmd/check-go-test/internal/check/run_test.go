// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
)

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{
			name: "go file",
			file: "file.go",
			want: true,
		},
		{
			name: "not go file",
			file: "file.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := check.IsGoFile(tc.file)

			if tc.want != got {
				t.Errorf("test.IsGoFile(%q) = %v; want %v", tc.file, got, tc.want)
			}
		})
	}
}

func TestCollectPackages(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{
			name: "no files",
		},
		{
			name: "one file",
			files: []string{
				"./package/file.go",
			},
			want: []string{"./package"},
		},
		{
			name: "adds dot prefix",
			files: []string{
				"package/file.go",
			},
			want: []string{"./package"},
		},
		{
			name: "removes duplicates",
			files: []string{
				"./package/file-a.go",
				"./package/file-b.go",
			},
			want: []string{"./package"},
		},
		{
			name: "keeps only go-packages",
			files: []string{
				"./package-a/file-a",
				"./package-b/file-b.go",
			},
			want: []string{"./package-b"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pkgs := check.CollectPackages(tc.files)

			if d := cmp.Diff(tc.want, pkgs); d != "" {
				t.Errorf("CollectPackages() mismatch (-want +got):\n%s", d)
				t.Logf("files:\n%v", tc.files)
			}
		})
	}
}
