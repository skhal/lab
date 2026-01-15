// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"testing"

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
