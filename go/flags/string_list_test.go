// Copyright 2025 Samvel Khalatyan. All rights reserved.

package flags_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/go/flags"
)

func TestStringList_Set(t *testing.T) {
	tests := []struct {
		name    string
		strings []string
		want    []string
	}{
		{
			name: "empty",
		},
		{
			name:    "one item",
			strings: []string{"one"},
			want:    []string{"one"},
		},
		{
			name:    "one item trim space right",
			strings: []string{"one "},
			want:    []string{"one"},
		},
		{
			name:    "one item trim space left",
			strings: []string{" one"},
			want:    []string{"one"},
		},
		{
			name:    "one item trim space",
			strings: []string{" one "},
			want:    []string{"one"},
		},
		{
			name:    "two items multiple flags",
			strings: []string{"one", "two"},
			want:    []string{"one", "two"},
		},
		{
			name:    "two items one flag with separator",
			strings: []string{"one,two"},
			want:    []string{"one", "two"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var flag flags.StringList

			for _, str := range tc.strings {
				err := flag.Set(str)
				if err != nil {
					t.Fatalf("flags.StringList().Set(%q) = %v; want no error", str, err)
				}
			}

			if diff := cmp.Diff(tc.want, flag.Get(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("flags.StringList(%q).Get() mismatch (-want, +got):\n%s", tc.strings, diff)
			}
		})
	}
}

func TestStringList_Get(t *testing.T) {
	tests := []struct {
		name    string
		strings []string
		want    []string
	}{
		{
			name: "empty",
		},
		{
			name:    "one item",
			strings: []string{"one"},
			want:    []string{"one"},
		},
		{
			name:    "two items",
			strings: []string{"one", "two"},
			want:    []string{"one", "two"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag := flags.StringList(tc.strings)

			got := flag.Get()

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("flags.StringList(%q).Get() mismatch (-want, +got):\n%s", tc.strings, diff)
			}
		})
	}
}

func TestStringList_String(t *testing.T) {
	tests := []struct {
		name    string
		strings []string
		want    string
	}{
		{
			name: "empty",
			want: "[]",
		},
		{
			name:    "one item",
			strings: []string{"one"},
			want:    "[one]",
		},
		{
			name:    "two items",
			strings: []string{"one", "two"},
			want:    "[one two]",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			flag := flags.StringList(tc.strings)

			got := flag.String()

			if got != tc.want {
				t.Errorf("flags.StringList(%s).String() = %q; want %q", tc.strings, got, tc.want)
			}
		})
	}
}
