// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check doc
package check

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var ErrNoDoc = errors.New("no documentation")

// Run verifies that every non-generated Go file has documentation attached to
// the exported declarations. It returns an error on the first failed file.
func Run(files []string) error {
	for _, f := range files {
		if IsTest(f) {
			continue
		}
		if err := Check(f); err != nil {
			return err
		}
	}
	return nil
}

// IsTest reports whether the file is a test file. A test file has _test.go
// suffix.
func IsTest(name string) bool {
	return strings.HasSuffix(name, "_test.go")
}

// Check verifies that non-generated Go file has documentation attached to the
// exported declarations. It returns an error if the check fails.
func Check(fname string) error {
	fset := token.NewFileSet()
	opts := parser.SkipObjectResolution | parser.ParseComments
	f, err := parser.ParseFile(fset, fname, nil, opts)
	if err != nil {
		return err
	}
	return CheckAST(fset, f)
}

// CheckAST scans top-level exported declarations for presence of documentation.
// It skips generated data.
func CheckAST(fset *token.FileSet, f *ast.File) error {
	if ast.IsGenerated(f) {
		return nil
	}
	return checkDecls(fset, f)
}

func checkDecls(fset *token.FileSet, f *ast.File) error {
	var ee []error
	for _, decl := range f.Decls {
		if err := checkDecl(fset, decl); err != nil {
			if !errors.Is(err, ErrNoDoc) {
				return err
			}
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func checkDecl(fset *token.FileSet, decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return checkFuncDecl(fset, d)
	}
	return nil
}

func checkFuncDecl(fset *token.FileSet, decl *ast.FuncDecl) error {
	if !decl.Name.IsExported() {
		return nil
	}
	if decl.Doc != nil {
		return nil
	}
	pos := fset.Position(decl.Pos())
	return fmt.Errorf("%s: func %s: %w", pos, decl.Name, ErrNoDoc)
}
