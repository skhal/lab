// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/skhal/lab/check/cmd/check-godoc/internal/check"
)

func TestIsTest(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "foo.go"},
		{name: "foo_test.go", want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := check.IsTest(tc.name)

			if got != tc.want {
				t.Errorf("check.IsTest(%q) got %v; want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestCheckAST(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "func exported no comment",
			code: `package test
func Test() {}`,
			want: check.ErrNoDoc,
		},
		{
			name: "func exported with comment",
			code: `package test
// Test comment
func Test() {}`,
		},
		{
			name: "func not exported no comment",
			code: `package test
func test() {}`,
		},
		{
			name: "func not exported with comment",
			code: `package test
// test comment
func test() {}`,
		},
		{
			name: "func exported no comment skip generated",
			code: `// Code generated ... DO NOT EDIT.
package test
func Test() {}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset, f := mustParse(t, tc.code)

			err := check.CheckAST(fset, f)

			if !errors.Is(err, tc.want) {
				t.Errorf("checksCheckAST() unexpected error %v; want %v", err, tc.want)
			}
		})
	}
}

func mustParse(t *testing.T, s string) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", s, parser.SkipObjectResolution|parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %s\ncode:\n%q\n", err, s)
	}
	return fset, f
}
