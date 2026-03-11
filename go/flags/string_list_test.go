// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/go/flags"
)

func ExampleStringList() {
	var tags []string
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	fs.Var(flags.NewStringList(&tags), "tag", "comma separated tags")
	err := fs.Parse([]string{"-tag", "1", "-tag", "2,3", "-tag", ",,4"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tags)
	// Output:
	// [1 2 3 4]
}

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
			var got []string
			sl := flags.NewStringList(&got)

			for _, str := range tc.strings {
				err := sl.Set(str)
				if err != nil {
					t.Fatalf("Set(%q) = %v; want no error", str, err)
				}
			}

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Set() mismatch (-want, +got):\n%s", diff)
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
		},
		{
			name:    "one item",
			strings: []string{"one"},
			want:    "one",
		},
		{
			name:    "two items",
			strings: []string{"one", "two"},
			want:    "one,two",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sl := flags.NewStringList(&tc.strings)

			got := sl.String()

			if got != tc.want {
				t.Errorf("String() = %q; want %q", got, tc.want)
			}
		})
	}
}
