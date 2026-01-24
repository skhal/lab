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
)

// ErrNoDoc represents missing documentation error.
var ErrNoDoc = errors.New("no documentation")

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
	case *ast.GenDecl:
		return checkGenDecl(fset, d)
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
	return newErrNoDoc(fset, decl.Name, kindFunc)
}

func checkGenDecl(fset *token.FileSet, decl *ast.GenDecl) error {
	switch decl.Tok {
	case token.TYPE:
		return checkGenDeclTypeSpec(fset, decl)
	case token.VAR:
		return checkGenDeclValueSpec(fset, decl)
	}
	return nil
}

func checkGenDeclValueSpec(fset *token.FileSet, decl *ast.GenDecl) error {
	// the comment might be in one of the two places:
	// ast.GenDecl.Doc: a group comment
	// ast.ValueSpec.Comment: a line comment
	var ee []error
	for _, spec := range decl.Specs {
		s := spec.(*ast.ValueSpec)
		if s.Comment != nil {
			continue
		}
		if decl.Doc != nil {
			continue
		}
		for _, n := range s.Names {
			if !n.IsExported() {
				continue
			}
			ee = append(ee, newErrNoDoc(fset, n, kindVar))
		}
	}
	return errors.Join(ee...)
}

func checkGenDeclTypeSpec(fset *token.FileSet, decl *ast.GenDecl) error {
	var ee []error
	for _, spec := range decl.Specs {
		s := spec.(*ast.TypeSpec)
		if !s.Name.IsExported() {
			continue
		}
		// A comment attached to a struct in a type group:
		//  type (
		//    // comment
		//    A struct {}
		//  }
		if s.Doc != nil {
			continue
		}
		// A comment attached to the type group, allow only one type inside:
		//  // comment
		//  type A struct {}
		if decl.Doc != nil && len(decl.Specs) == 1 {
			continue
		}
		ee = append(ee, newErrNoDoc(fset, s.Name, kindType))
	}
	return errors.Join(ee...)
}

type kind int

//go:generate stringer -type=kind -linecomment
const (
	_ kind = iota
	// keep-sorted start
	kindFunc // func
	kindType // type
	kindVar  // var
	// keep-sorted end
)

func newErrNoDoc(fset *token.FileSet, id *ast.Ident, k kind) error {
	pos := fset.Position(id.Pos())
	return fmt.Errorf("%s: %s %s: %w", pos, k, id.Name, ErrNoDoc)
}
