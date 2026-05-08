// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nosubmit_test

import (
	"errors"
	"os"
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

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr error
	}{
		{
			name:    "missing file",
			file:    "testdata/does_not_exist",
			wantErr: os.ErrNotExist,
		},
		{
			name:    "missing file",
			file:    "testdata/has_nosubmit.txt",
			wantErr: nosubmit.ErrCheck,
		},
		{
			name: "missing file",
			file: "testdata/has_no_nosubmit.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := nosubmit.Run(tc.file)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error %v; want %v", err, tc.wantErr)
			}
		})
	}
}
