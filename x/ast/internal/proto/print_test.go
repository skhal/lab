// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto_test

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	pbast "github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/x/ast/internal/proto"
)

var update = flag.Bool("update", false, "update golden files")

func TestFprint(t *testing.T) {
	tests := []struct {
		name   string
		proto  string
		golden tests.GoldenFile
	}{
		{
			name:   "empty",
			golden: "testdata/empty.txt",
		},
		{
			name: "edition",
			proto: `
edition = "2024";
package test;
`,
			golden: "testdata/edition.txt",
		},
		{
			name: "syntax",
			proto: `
syntax = "proto3";
package test;
`,
			golden: "testdata/syntax.txt",
		},
		{
			name: "message",
			proto: `
edition = "2024";
package test;
message Test {
	string one = 1;
}
`,
			golden: "testdata/message.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file := mustParse(t, tc.proto)
			var buf bytes.Buffer

			proto.Fprint(&buf, file)

			if *update {
				tc.golden.Write(t, buf.String())
			}
			if diff := tc.golden.Diff(t, buf.String()); diff != "" {
				t.Errorf("proto.Fprint() unexpected diff (-want,+got):\n%s", diff)
				t.Log(tc.proto)
			}
		})
	}
}

func mustParse(t *testing.T, s string) *pbast.FileNode {
	t.Helper()
	r := strings.NewReader(s)
	f, err := parser.Parse("test", r, reporter.NewHandler(nil))
	if err != nil {
		t.Fatal(err)
	}
	return f
}
