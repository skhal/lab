// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palindrome_test

import (
	"testing"

	"github.com/skhal/lab/iq/3/palindrome"
)

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "one letter",
			str:  "a",
			want: true,
		},
		{
			name: "one digit",
			str:  "1",
			want: true,
		},
		{
			name: "non-alphanumic letter",
			str:  ".",
			want: true,
		},
		{
			name: "two letters",
			str:  "aa",
			want: true,
		},
		{
			name: "two digits",
			str:  "11",
			want: true,
		},
		{
			name: "two distinct letters",
			str:  "ab",
		},
		{
			name: "letter and digit",
			str:  "a1",
		},
		{
			name: "letter and non-alphanumeric",
			str:  "a.",
			want: true,
		},
		{
			name: "digit and non-alphanumeric",
			str:  "1.",
			want: true,
		},
		{
			name: "non-alphanumeric and letter",
			str:  ".a",
			want: true,
		},
		{
			name: "non-alphanumeric and digit",
			str:  ".1",
			want: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := palindrome.Is(tc.str)
			if got != tc.want {
				t.Errorf("palindrome.Is(%q) = %v, want %v", tc.str, got, tc.want)
			}
		})
	}
}
