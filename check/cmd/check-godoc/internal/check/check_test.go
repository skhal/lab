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

func TestCheckAST_func(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "exported no comment",
			code: `package test
func Test() {}`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported with comment",
			code: `package test
// Test comment
func Test() {}`,
		},
		{
			name: "not exported no comment",
			code: `package test
func test() {}`,
		},
		{
			name: "not exported with comment",
			code: `package test
// test comment
func test() {}`,
		},
		{
			name: "exported no comment skip generated",
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
				t.Log(tc.code)
			}
		})
	}
}

func TestCheckAST_var(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "exported no comment",
			code: `package test
var A = 1`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported comment",
			code: `package test
// test comment
var A = 1`,
		},
		{
			name: "exported line comment",
			code: `package test
var A = 1 // test comment`,
		},
		{
			name: "not exported no comment",
			code: `package test
var a = 1`,
		},
		{
			name: "exported multiple no comment",
			code: `package test
var A, B = 1, 1`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported multiple comment",
			code: `package test
// test comment
var A, B = 1, 1`,
		},
		{
			name: "exported multiple line comment",
			code: `package test
var A, B = 1, 1 // test comment`,
		},
		{
			name: "exported multiple mixed no comment",
			code: `package test
var A, b = 1, 1`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported multiple mixed comment",
			code: `package test
// test comment
var A, b = 1, 1`,
		},
		{
			name: "exported multiple mixed line comment",
			code: `package test
var A, b = 1, 1 // test comment`,
		},
		{
			name: "exported group no comment",
			code: `package test
var (
	A = 1
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported group comment",
			code: `package test
// test comment
var (
	A = 1
)`,
		},
		{
			name: "exported group line comment",
			code: `package test
var (
	A = 1 // test comment
)`,
		},
		{
			name: "exported group multiple no comment",
			code: `package test
var (
	A = 1
	B = 1
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported group multiple comment",
			code: `package test
// test comment
var (
	A = 1
	B = 1
)`,
		},
		{
			name: "exported group multiple line comment",
			code: `package test
var (
	A = 1 // test comment
	B = 1 // test comment
)`,
		},
		{
			name: "exported group multiple some line comment",
			code: `package test
var (
	A = 1 // test comment
	B = 1
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported group mixed comment",
			code: `package test
// test comment
var (
	A = 1
	b = 1
)`,
		},
		{
			name: "exported group mixed line comment",
			code: `package test
var (
	A = 1 // test comment
	b = 1 // test comment
)`,
		},
		{
			name: "exported group mixed some line comment",
			code: `package test
var (
	A = 1 // test comment
	b = 1
)`,
		},
		{
			name: "exported group mixed no line comment",
			code: `package test
var (
	A = 1
	b = 1 // test comment
)`,
			want: check.ErrNoDoc,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset, f := mustParse(t, tc.code)

			err := check.CheckAST(fset, f)

			if !errors.Is(err, tc.want) {
				t.Errorf("checksCheckAST() unexpected error %v; want %v", err, tc.want)
				t.Log(tc.code)
			}
		})
	}
}

func TestCheckAST_struct(t *testing.T) {
	tests := []struct {
		name string
		code string
		want error
	}{
		{
			name: "exported no comment",
			code: `package test
type A struct {}`,
			want: check.ErrNoDoc,
		},
		{
			name: "exported comment",
			code: `package test
// test comment
type A struct {}`,
		},
		{
			name: "exported line comment",
			code: `package test
type A struct {} // test comment`,
			want: check.ErrNoDoc,
		},
		{
			name: "not exported no comment",
			code: `package test
type a struct {}`,
		},
		{
			name: "group one exported no comment",
			code: `package test
type (
	A struct {}
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "group one exported comment",
			code: `package test
// test comment
type (
  A struct {}
)`,
		},
		{
			name: "group one exported line comment",
			code: `package test
type (
  A struct {} // test comment
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "group one not exported no comment",
			code: `package test
type (
  a struct {}
)`,
		},
		{
			name: "group two exported no comment",
			code: `package test
type (
	A struct {}
	B struct {}
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "group two exported group comment",
			code: `package test
// test comment
type (
  A struct {}
	B struct {}
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "group two exported line comment",
			code: `package test
type (
  A struct {} // test comment
  B struct {} // test comment
)`,
			want: check.ErrNoDoc,
		},
		{
			name: "group two exported",
			code: `package test
type (
	// A test
  A struct {}
	// B test
  B struct {}
)`,
		},
		{
			name: "group two not exported no comment",
			code: `package test
type (
  a struct {}
  b struct {}
)`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset, f := mustParse(t, tc.code)

			err := check.CheckAST(fset, f)

			if !errors.Is(err, tc.want) {
				t.Errorf("checksCheckAST() unexpected error %v; want %v", err, tc.want)
				t.Log(tc.code)
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
