// Copyright 2025 Samvel Khalatyan. All rights reserved.

package anagram_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/string/anagram"
)

func TestFindAll(t *testing.T) {
	tests := []struct {
		name string
		s    string
		t    string
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "single match",
			s:    "ab",
			t:    "ab",
			want: []string{"ab"},
		},
		{
			name: "single match with prefix",
			s:    "aab",
			t:    "ab",
			want: []string{"ab"},
		},
		{
			name: "single match with suffix",
			s:    "abb",
			t:    "ab",
			want: []string{"ab"},
		},
		{
			name: "two matches",
			s:    "aba",
			t:    "ab",
			want: []string{"ab", "ba"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := anagram.FindAll(tc.s, tc.t)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("anagram.FindAll(%q, %q) mismatch (-want, +got):\n%s", tc.s, tc.t, diff)
			}
		})
	}
}
