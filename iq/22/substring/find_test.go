// Copyright 2025 Samvel Khalatyan. All rights reserved.

package substring_test

import (
	"testing"

	"github.com/skhal/lab/iq/22/substring"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "empty",
		},
		{
			name: "len one",
			s: "a",
			want: "a",
		},
		{
			name: "len two",
			s: "ab",
			want: "ab",
		},
		{
			name: "len two same char",
			s: "aa",
			want: "a",
		},
		{
			name: "len three",
			s: "abc",
			want: "abc",
		},
		{
			name: "len three same char",
			s: "aaa",
			want: "a",
		},
		{
			name: "len three first two chars",
			s: "aba",
			want: "ab",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := substring.Find(tc.s)

			if got != tc.want {
				t.Errorf("substring.Find(%q) = %q; want %q", tc.s, got, tc.want)
			}
		})
	}
}
