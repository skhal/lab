// Copyright 2025 Samvel Khalatyan. All rights reserved.

package test_test

import (
	"testing"

	"github.com/skhal/lab/check/cmd/gotest/internal/test"
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
			got := test.IsGoFile(tc.file)

			if tc.want != got {
				t.Errorf("test.IsGoFile(%q) = %v; want %v", tc.file, got, tc.want)
			}
		})
	}
}
