// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dups_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/skhal/lab/iq/38/dups"
)

func ExampleRemove() {
	s := "abccbaad"
	got := dups.Remove(s)
	fmt.Printf("%q .. %q\n", s, got)
	// Output:
	// "abccbaad" .. "ad"
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
		{s: "aaa", want: "a"},
		{s: "aab", want: "b"},
		{s: "baa", want: "b"},
		{s: "aba", want: "aba"},
		{s: "aaaa", want: ""},
		{s: "aaab", want: "ab"},
		{s: "aaba", want: "ba"},
		{s: "abaa", want: "ab"},
		{s: "baaa", want: "ba"},
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
