// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/skhal/lab/check/cmd/check-pbdoc/internal/check"
)

func TestCheckFile_message(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name: "message with comment",
			s: `
edition = "2024";
package test;
// Test comment
message Test {}
`,
		},
		{
			name: "message with comment wrong prefix",
			s: `
edition = "2024";
package test;
// a comment
message Test {}
`,
			wantErr: check.ErrCommentPrefix,
		},
		{
			name: "message no comment",
			s: `
edition = "2024";
package test;
message Test {}
`,
			wantErr: check.ErrNoComment,
		},
		{
			name: "message field with comment",
			s: `
edition = "2024";
package test;
// Foo comment
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
// Foo comment
message Foo {
	string test = 1;
}
`,
			wantErr: check.ErrNoComment,
		},
		{
			name: "nested message with comment",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	// Test comment
	message Test {}
}
`,
		},
		{
			name: "nested message with comment wrong prefix",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	// A comment
	message Test {}
}
`,
			wantErr: check.ErrCommentPrefix,
		},
		{
			name: "nested message no comment",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	message Test {}
}
`,
			wantErr: check.ErrNoComment,
		},
		{
			name: "nested message field with comment",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	// Bar comment
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
// Foo comment
message Foo {
	// Bar comment
	message Bar {
		string test = 1;
	}
}
`,
			wantErr: check.ErrNoComment,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := check.CheckFile("test.proto", strings.NewReader(tc.s))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("CheckFile() unexpected error %v, want %v", err, tc.wantErr)
				t.Log(tc.s)
			}
		})
	}
}

func TestCheckFile_enum(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr error
	}{
		{
			name: "enum with comment",
			s: `
edition = "2024";
package test;
// Test comment
enum Test {}
`,
		},
		{
			name: "enum with comment wrong prefix",
			s: `
edition = "2024";
package test;
// a comment
enum Test {}
`,
			wantErr: check.ErrCommentPrefix,
		},
		{
			name: "enum no comment",
			s: `
edition = "2024";
package test;
enum Test {}
`,
			wantErr: check.ErrNoComment,
		},
		{
			name: "enumerator with comment",
			s: `
edition = "2024";
package test;
// Foo comment
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
// Foo comment
enum Foo {
	FOO_TEST = 1;
}
`,
			wantErr: check.ErrNoComment,
		},
		{
			name: "nested enum with comment",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	// Test comment
	enum Test {}
}
`,
		},
		{
			name: "nested enum with comment wrong prefix",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	// a comment
	enum Test {}
}
`,
			wantErr: check.ErrCommentPrefix,
		},
		{
			name: "nested enum no comment",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	enum Test {}
}
`,
			wantErr: check.ErrNoComment,
		},
		{
			name: "nested enum enumerator with comment",
			s: `
edition = "2024";
package test;
// Foo comment
message Foo {
	// Bar comment
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
// Foo comment
message Foo {
	// Bar comment
	enum Bar {
		BAR_TEST = 1;
	}
}
`,
			wantErr: check.ErrNoComment,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := check.CheckFile("test.proto", strings.NewReader(tc.s))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("CheckFile() unexpected error %v, want %v", err, tc.wantErr)
				t.Log(tc.s)
			}
		})
	}
}
