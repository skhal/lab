// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexiseq_test

import (
	"testing"

	"github.com/skhal/lab/iq/6/lexiseq"
)

func TestNext(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "empty",
		},
		{
			name: "one letter",
			s:    "a",
			want: "a",
		},
		// Two letters
		{
			name: "two letters same",
			s:    "aa",
			want: "aa",
		},
		{
			name: "two letters ascending",
			s:    "ab",
			want: "ba",
		},
		{
			name: "two letters descending",
			s:    "ba",
			want: "ab",
		},
		// Three letters
		{
			name: "three letters same",
			s:    "aaa",
			want: "aaa",
		},
		{
			name: "three letters first two same",
			s:    "aab",
			want: "aba",
		},
		{
			name: "three letters first and last same",
			s:    "aba",
			want: "baa",
		},
		{
			name: "three letters second and last same",
			s:    "baa",
			want: "aab",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := lexiseq.Next(tc.s)
			if tc.want != got {
				t.Errorf("lexiseq.Next(%q) = %q, want %q", tc.s, got, tc.want)
			}
		})
	}
}
