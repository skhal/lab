// Copyright 2025 Samvel Khalatyan. All rights reserved.

package dups_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/skhal/lab/iq/37/dups"
)

func ExampleRemove() {
	s := "abccbad"
	got := dups.Remove(s)
	fmt.Printf("%q .. %q\n", s, got)
	// Output:
	// "abccbad" .. "d"
}

func TestRemove(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{}, // empty
		{s: "a", want: "a"},
		{s: "aa", want: ""},
		{s: "ab", want: "ab"},
		{s: "aab", want: "b"},
		{s: "baa", want: "b"},
		{s: "aba", want: "aba"},
		{s: "aaaa", want: ""},
		{s: "aaab", want: "b"},
		{s: "aaba", want: "ba"},
		{s: "abaa", want: "ab"},
		{s: "baaa", want: "b"},
		{s: "abba", want: ""},
		{s: "abbc", want: "ac"},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := dups.Remove(tc.s)

			if tc.want != got {
				t.Errorf("dups.Remove(%q) = %q; want %q", tc.s, got, tc.want)
			}
		})
	}
}
