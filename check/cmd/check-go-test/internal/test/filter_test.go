// Copyright 2025 Samvel Khalatyan. All rights reserved.

package test_test

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestFilterFunc(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		fn    func(string) bool
		want  []string
	}{
		{
			name: "empty",
			fn:   func(string) bool { return true },
		},
		{
			name:  "pass all",
			items: []string{"a", "b"},
			fn:    func(string) bool { return true },
			want:  []string{"a", "b"},
		},
		{
			name:  "pass none",
			items: []string{"a", "b"},
			fn:    func(string) bool { return false },
		},
		{
			name:  "pass some",
			items: []string{"a", "b", "c"},
			fn:    func(s string) bool { return s == "b" },
			want:  []string{"b"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Collect(test.FilterFunc(slices.Values(tc.items), tc.fn))

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("test.FilterFunc() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestPaths(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{
			name: "empty",
		},
		{
			name:  "base only",
			files: []string{"file.txt"},
			want:  []string{""},
		},
		{
			name:  "path and base",
			files: []string{"path/file.txt"},
			want:  []string{"path/"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Collect(test.Paths(slices.Values(tc.files)))

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("test.Paths() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  []string
	}{
		{
			name: "empty",
		},
		{
			name:  "all different",
			items: []string{"a", "b"},
			want:  []string{"a", "b"},
		},
		{
			name:  "dups",
			items: []string{"a", "b", "a"},
			want:  []string{"a", "b"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Collect(test.Unique(slices.Values(tc.items)))

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("test.Unique() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
