// Copyright 2025 Samvel Khalatyan. All rights reserved.

package substring_test

import (
	"testing"

	"github.com/skhal/lab/iq/23/substring"
)

type testCase struct {
	name string
	s    string
	n    int
	want string
}

func TestFind(t *testing.T) {
	tests := []struct {
		group string
		tests []testCase
	}{
		{
			group: "n=0",
			tests: []testCase{
				{
					name: "empty",
				},
				{
					name: "one letter",
					s:    "a",
					want: "a",
				},
				{
					name: "two letters",
					s:    "ab",
					want: "a",
				},
				{
					name: "two letters same",
					s:    "aa",
					want: "aa",
				},
				{
					name: "three letters",
					s:    "abc",
					want: "a",
				},
				{
					name: "three letters first two same",
					s:    "aac",
					want: "aa",
				},
				{
					name: "three letters first and last same",
					s:    "aba",
					want: "a",
				},
				{
					name: "three letters last two same",
					s:    "abb",
					want: "bb",
				},
				{
					name: "three letters same",
					s:    "aaa",
					want: "aaa",
				},
			},
		},
		{
			group: "n=1",
			tests: []testCase{
				{
					name: "empty",
					n:    1,
				},
				{
					name: "one letter",
					s:    "a",
					n:    1,
					want: "a",
				},
				{
					name: "two letters",
					s:    "ab",
					n:    1,
					want: "ab",
				},
				{
					name: "two letters same",
					s:    "aa",
					n:    1,
					want: "aa",
				},
				{
					name: "three letters",
					s:    "abc",
					n:    1,
					want: "ab",
				},
				{
					name: "three letters first two same",
					s:    "aac",
					n:    1,
					want: "aac",
				},
				{
					name: "three letters first and last same",
					s:    "aba",
					n:    1,
					want: "aba",
				},
				{
					name: "three letters last two same",
					s:    "abb",
					n:    1,
					want: "abb",
				},
				{
					name: "three letters same",
					s:    "aaa",
					n:    1,
					want: "aaa",
				},
			},
		},
	}
	for _, gtc := range tests {
		gtc := gtc
		t.Run(gtc.group, func(t *testing.T) {
			for _, tc := range gtc.tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					got := substring.Find(tc.s, tc.n)

					if got != tc.want {
						t.Errorf("substring.Find(%q) = %q; want %q", tc.s, got, tc.want)
					}
				})
			}
		})
	}
}
