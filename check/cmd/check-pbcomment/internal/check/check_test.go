// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"strings"
	"testing"

	"github.com/skhal/lab/check/cmd/check-pbcomment/internal/check"
)

func TestCheckFile(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{
			name: "empty",
		},
		{
			name: "message with comment",
			s: `
edition = "2024";
package test;
// test comment
message Test {}
`,
		},
		{
			name: "message no comment",
			s: `
edition = "2024";
package test;
message Test {}
`,
			wantErr: true,
		},
		{
			name: "enum with comment",
			s: `
edition = "2024";
package test;
// test comment
enum Test {}
`,
		},
		{
			name: "enum no comment",
			s: `
edition = "2024";
package test;
enum Test {}
`,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := check.CheckFile("test.proto", strings.NewReader(tc.s))

			if err == nil {
				if tc.wantErr {
					t.Error("CheckFile() want error")
					t.Log(tc.s)
				}
			} else if !tc.wantErr {
				t.Errorf("CheckFile() unexpected error %v", err)
				t.Log(tc.s)
			}
		})
	}
}
