// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"strings"
	"testing"

	"github.com/skhal/lab/check/cmd/check-pbdoc/internal/check"
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
			name: "message field with comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// test comment
	string test = 1;
}
`,
		},
		{
			name: "message field no comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	string test = 1;
}
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
		{
			name: "enumerator with comment",
			s: `
edition = "2024";
package test;
// foo comment
enum Foo {
	// test comment
	FOO_TEST = 1;
}
`,
		},
		{
			name: "enumerator no comment",
			s: `
edition = "2024";
package test;
// foo comment
enum Foo {
	FOO_TEST = 1;
}
`,
			wantErr: true,
		},
		{
			name: "nested message with comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// test comment
	message Test {}
}
`,
		},
		{
			name: "nested message no comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	message Test {}
}
`,
			wantErr: true,
		},
		{
			name: "nested message field with comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// bar comment
	message Bar {
		// test comment
		string test = 1;
	}
}
`,
		},
		{
			name: "nested message field no comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// bar comment
	message Bar {
		string test = 1;
	}
}
`,
			wantErr: true,
		},
		{
			name: "nested enum with comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// test comment
	enum Test {}
}
`,
		},
		{
			name: "nested enum no comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	enum Test {}
}
`,
			wantErr: true,
		},
		{
			name: "nested enum enumerator with comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// bar comment
	enum Bar {
		// test comment
		BAR_TEST = 1;
	}
}
`,
		},
		{
			name: "nested enum enumerator no comment",
			s: `
edition = "2024";
package test;
// foo comment
message Foo {
	// bar comment
	enum Bar {
		BAR_TEST = 1;
	}
}
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
