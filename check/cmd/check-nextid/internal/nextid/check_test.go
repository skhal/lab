// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nextid_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/skhal/lab/check/cmd/check-nextid/internal/nextid"
)

func TestCheckFileNode_message(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "empty message no next-id",
			code: `
edition = "2024";
package test;
message Test {}
`,
		},
		{
			name: "empty message valid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 1
message Test {}
`,
		},
		{
			name: "empty message invalid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 2
message Test {}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "not empty message no next-id",
			code: `
edition = "2024";
package test;
message Test {
	int foo = 1;
}
`,
		},
		{
			name: "not empty message valid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 2
message Test {
	int foo = 1;
}
`,
		},
		{
			name: "not empty message invalid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 3
message Test {
	int foo = 1;
}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "message with one reserved",
			code: `
edition = "2024";
package test;
// Next ID: 2
message Test {
	reserved 1;
}
`,
		},
		{
			name: "message with few reserved",
			code: `
edition = "2024";
package test;
// Next ID: 4
message Test {
	reserved 1,3;
}
`,
		},
		{
			name: "message with range reserved",
			code: `
edition = "2024";
package test;
// Next ID: 4
message Test {
	reserved 1 to 3;
}
`,
		},
		{
			name: "message with range mix",
			code: `
edition = "2024";
package test;
// Next ID: 6
message Test {
	reserved 1, 2 to 3, 5;
}
`,
		},
		{
			name: "message with field and range",
			code: `
edition = "2024";
package test;
// Next ID: 4
message Test {
	int foo = 2;
	reserved 1, 3;
}
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := nextid.CheckFileNode(mustParse(t, tc.code))

			if !errors.Is(err, tc.want) {
				t.Errorf("nextid.CheckFileNode() unexpected error %v; want %v", err, tc.want)
				t.Log(tc.code)
			}
		})
	}
}

func TestCheckFileNode_enum(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "empty enum no next-id",
			code: `
edition = "2024";
package test;
enum Test {}
`,
		},
		{
			name: "empty enum valid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 1
enum Test {}
`,
		},
		{
			name: "empty enum invalid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 2
enum Test {}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "not empty enum no next-id",
			code: `
edition = "2024";
package test;
enum Test {
	TEST_FOO = 1;
}
`,
		},
		{
			name: "not empty enum valid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 2
enum Test {
	TEST_FOO = 1;
}
`,
		},
		{
			name: "not empty enum invalid next-id",
			code: `
edition = "2024";
package test;
// Next ID: 3
enum Test {
	TEST_FOO = 1;
}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "enum with one reserved",
			code: `
edition = "2024";
package test;
// Next ID: 2
enum Test {
	reserved 1;
}
`,
		},
		{
			name: "enum with few reserved",
			code: `
edition = "2024";
package test;
// Next ID: 4
enum Test {
	reserved 1,3;
}
`,
		},
		{
			name: "enum with range reserved",
			code: `
edition = "2024";
package test;
// Next ID: 4
enum Test {
	reserved 1 to 3;
}
`,
		},
		{
			name: "enum with range mix",
			code: `
edition = "2024";
package test;
// Next ID: 6
enum Test {
	reserved 1, 2 to 3, 5;
}
`,
		},
		{
			name: "enum with field and range",
			code: `
edition = "2024";
package test;
// Next ID: 4
enum Test {
	TEST_FOO = 1;
	reserved 1, 3;
}
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := nextid.CheckFileNode(mustParse(t, tc.code))

			if !errors.Is(err, tc.want) {
				t.Errorf("nextid.CheckFileNode() unexpected error %v; want %v", err, tc.want)
				t.Log(tc.code)
			}
		})
	}
}

func TestCheckFileNode_submessage(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "empty message no next-id",
			code: `
edition = "2024";
package test;
message Test {
	message Child {}
}
`,
		},
		{
			name: "empty message valid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 1
	message Child {}
}
`,
		},
		{
			name: "empty message invalid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 2
	message Child {}
}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "not empty message no next-id",
			code: `
edition = "2024";
package test;
message Test {
	message Child {
		int foo = 1;
	}
}
`,
		},
		{
			name: "not empty message valid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 2
	message Child {
		int foo = 1;
	}
}
`,
		},
		{
			name: "not empty message invalid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 3
	message Child {
		int foo = 1;
	}
}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "message with one reserved",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 2
	message Child {
		reserved 1;
	}
}
`,
		},
		{
			name: "message with few reserved",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 4
	message Child {
		reserved 1,3;
	}
}
`,
		},
		{
			name: "message with range reserved",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 4
	message Child {
		reserved 1 to 3;
	}
}
`,
		},
		{
			name: "message with range mix",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 6
	message Child {
		reserved 1, 2 to 3, 5;
	}
}
`,
		},
		{
			name: "message with field and range",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 4
	message Child {
		int foo = 2;
		reserved 1, 3;
	}
}
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := nextid.CheckFileNode(mustParse(t, tc.code))

			if !errors.Is(err, tc.want) {
				t.Errorf("nextid.CheckFileNode() unexpected error %v; want %v", err, tc.want)
				t.Log(tc.code)
			}
		})
	}
}

func TestCheckFileNode_subenum(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "empty enum no next-id",
			code: `
edition = "2024";
package test;
message Test {
	enum Child {}
}
`,
		},
		{
			name: "empty enum valid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 1
	enum Child {}
}
`,
		},
		{
			name: "empty enum invalid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 2
	enum Child {}
}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "not empty enum no next-id",
			code: `
edition = "2024";
package test;
message Test {
	enum Child {
		CHILD_FOO = 1;
	}
}
`,
		},
		{
			name: "not empty enum valid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 2
	enum Child {
		CHILD_FOO = 1;
	}
}
`,
		},
		{
			name: "not empty enum invalid next-id",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 3
	enum Child {
		CHILD_FOO = 1;
	}
}
`,
			want: nextid.ErrNextID,
		},
		{
			name: "enum with one reserved",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 2
	enum Child {
		reserved 1;
	}
}
`,
		},
		{
			name: "enum with few reserved",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 4
	enum Child {
		reserved 1,3;
	}
}
`,
		},
		{
			name: "enum with range reserved",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 4
	enum Child {
		reserved 1 to 3;
	}
}
`,
		},
		{
			name: "enum with range mix",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 6
	enum Child {
		reserved 1, 2 to 3, 5;
	}
}
`,
		},
		{
			name: "enum with field and range",
			code: `
edition = "2024";
package test;
message Test {
	// Next ID: 4
	enum Child {
		CHILD_FOO = 1;
		reserved 1, 3;
	}
}
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := nextid.CheckFileNode(mustParse(t, tc.code))

			if !errors.Is(err, tc.want) {
				t.Errorf("nextid.CheckFileNode() unexpected error %v; want %v", err, tc.want)
				t.Log(tc.code)
			}
		})
	}
}

func mustParse(t *testing.T, s string) *ast.FileNode {
	t.Helper()
	fn, err := parser.Parse("string", strings.NewReader(s), reporter.NewHandler(nil))
	if err != nil {
		t.Fatal(err)
	}
	return fn
}

func TestParseNextID(t *testing.T) {
	type want struct {
		id  uint64
		ok  bool
		err error
	}
	tests := []struct {
		name    string
		comment string
		want    want
	}{
		{
			name:    "empty",
			comment: ``,
		},
		{
			name:    "valid",
			comment: `// Next ID: 2`,
			want: want{
				id: 2,
				ok: true,
			},
		},
		{
			name:    "valid mixed case",
			comment: `// nExT id: 2`,
			want: want{
				id: 2,
				ok: true,
			},
		},
		{
			name:    "negative is invalid",
			comment: `// Next ID: -2`,
			want: want{
				err: nextid.ErrNextID,
			},
		},
		{
			name:    "zero is invalid",
			comment: `// Next ID: 0`,
			want: want{
				err: nextid.ErrNextID,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, ok, err := nextid.ParseNextID(tc.comment)

			if id != tc.want.id {
				t.Errorf("nextid.ParseNextID(%q) = %d, _, _; want %d", tc.comment, id, tc.want.id)
			}
			if ok != tc.want.ok {
				t.Errorf("nextid.ParseNextID(%q) = _, %v, _; want %v", tc.comment, ok, tc.want.ok)
			}
			if !errors.Is(err, tc.want.err) {
				t.Errorf("nextid.ParseNextID(%q) = _, _, %v; want %d", tc.comment, err, tc.want.err)
			}
		})
	}
}
