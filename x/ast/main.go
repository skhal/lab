// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ast pases a Go file and prints Abstract Syntax Tree.
//
// Usage:
//
//	ast file
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	pbparser "github.com/bufbuild/protocompile/parser"
	pbreporter "github.com/bufbuild/protocompile/reporter"
	"github.com/skhal/lab/x/ast/internal/proto"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s file\n", filepath.Base(os.Args[0]))
	}
	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(name string) error {
	ok, err := isDir(name)
	switch {
	case err != nil:
		return err
	case ok:
		return runOnDir(name)
	default:
		return runOnFile(name)
	}
}

func isDir(name string) (bool, error) {
	m, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return m.IsDir(), nil
}

func runOnDir(name string) error {
	fset := token.NewFileSet()
	opts := parser.ParseComments | parser.PackageClauseOnly
	noTests := func(m fs.FileInfo) bool {
		return !strings.HasSuffix(m.Name(), "_test.go")
	}
	f, err := parser.ParseDir(fset, name, noTests, opts)
	if err != nil {
		return err
	}
	ast.Print(fset, f)
	return nil
}

func runOnFile(name string) error {
	switch filepath.Ext(name) {
	case ".go":
		return runOnGoFile(name)
	case ".proto":
		return runOnProtoFile(name)
	default:
		return fmt.Errorf("unsupported file %s", name)
	}
}

func runOnGoFile(name string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Print(fset, f)
	return nil
}

func runOnProtoFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	file, err := pbparser.Parse(name, f, pbreporter.NewHandler(nil))
	if err != nil {
		return err
	}
	proto.Print(file)
	return nil
}
