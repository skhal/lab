// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nosubmit_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/check/cmd/check-nosubmit/internal/nosubmit"
)

func ExampleRun() {
	readFileFn := func(f string) ([]byte, error) {
		if f != "foo.txt" {
			return nil, fmt.Errorf("error opening file %s", f)
		}
		s := `
// DO NOT SUBMIT: work in progress
		`
		return []byte(s), nil
	}
	cfg := &nosubmit.Config{
		ReadFileFn: readFileFn,
	}
	if err := nosubmit.Run(cfg, "foo.txt"); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// check error
}

func TestHasNoSubmit(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{name: "empty"},
		{
			name: "pass",
			data: `
test data
`,
		},
		{
			name: "nosubmit",
			data: `
test data
// DO NOT SUBMIT
`,
			want: true,
		},
		{
			name: "nosubmit with comment",
			data: `
test data
// DO NOT SUBMIT: description
`,
			want: true,
		},
		{
			name: "raw nosubmit",
			data: `
test data
DO NOT SUBMIT
`,
			want: true,
		},
		{
			name: "raw nosubmit with comment",
			data: `
test data
DO NOT SUBMIT: description
`,
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := nosubmit.Check([]byte(tc.data))

			if tc.want != got {
				t.Errorf("nosubmit.Check(...) = %v; want %v", got, tc.want)
				t.Logf("data:\n%s", tc.data)
			}
		})
	}
}
