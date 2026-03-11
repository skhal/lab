// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nosubmit_test

import (
	"testing"

	"github.com/skhal/lab/check/cmd/check-nosubmit/internal/nosubmit"
)

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
				t.Errorf("Check() = %v; want %v", got, tc.want)
				t.Logf("data:\n%s", tc.data)
			}
		})
	}
}
