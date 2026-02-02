// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package proto implements functions to print Protobuf AST.
package proto

import (
	"fmt"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/bufbuild/protocompile/ast"
)

// Print dumps the file node fn to standard output.
func Print(fn *ast.FileNode) {
	fprint(os.Stdout, 0, reflect.ValueOf(fn))
}

// Fprint dumps the file node fn to the writer w.
func Fprint(w io.Writer, fn *ast.FileNode) {
	fprint(w, 0, reflect.ValueOf(fn))
}

func fprint(w io.Writer, lvl int, val reflect.Value) {
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		t := val.Type()
		fmt.Fprintf(w, "%s {\n", t)
		for idx := range val.Len() {
			fmt.Fprintf(w, "%s%d: ", genPrefix(lvl+1), idx)
			fprint(w, lvl+1, val.Index(idx))
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "%s}", genPrefix(lvl))
	case reflect.Interface:
		fprint(w, lvl, val.Elem())
	case reflect.Pointer:
		fmt.Fprint(w, "*")
		fprint(w, lvl, val.Elem())
	case reflect.String:
		fmt.Fprintf(w, "%q", val)
	case reflect.Struct:
		t := val.Type()
		fmt.Fprintf(w, "%s {\n", t)
		for idx := range t.NumField() {
			name := t.Field(idx).Name
			if !token.IsExported(name) {
				continue
			}
			f := val.Field(idx)
			if skip(f) {
				continue
			}
			fmt.Fprintf(w, "%s%s: ", genPrefix(lvl+1), name)
			fprint(w, lvl+1, f)
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "%s}", genPrefix(lvl))
	default:
		fmt.Fprint(w, val)
	}
}

func genPrefix(lvl int) string {
	return strings.Repeat(". ", lvl)
}

func skip(v reflect.Value) bool {
	if !v.IsValid() || v.IsZero() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}
