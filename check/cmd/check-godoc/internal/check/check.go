// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"maps"
	"slices"
	"strings"
)

// ErrNoDoc represents missing documentation error.
var ErrNoDoc = errors.New("no documentation")

// ErrNotPackage means that Go failed to parse a path.
var ErrNotPackage = errors.New("not a package")

// ErrMultiPackage means there are multiple packages at a single path.
var ErrMultiPackage = errors.New("multiple packages")

// ErrMultiDoc means there are multiple package documentations.
var ErrMultiDoc = errors.New("multiple documentation")

// CheckFile verifies that non-generated Go file has documentation attached to
// the exported declarations. It returns an error if the check fails.
func CheckFile(fname string) error {
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
	case token.CONST:
		return checkGenDeclConstValueSpec(fset, decl)
	case token.TYPE:
		return checkGenDeclTypeSpec(fset, decl)
	case token.VAR:
		return checkGenDeclVarValueSpec(fset, decl)
	}
	return nil
}
func checkGenDeclConstValueSpec(fset *token.FileSet, decl *ast.GenDecl) error {
	var ee []error
	for _, spec := range decl.Specs {
		s := spec.(*ast.ValueSpec)
		for _, n := range s.Names {
			if !n.IsExported() {
				continue
			}
			if s.Doc != nil && len(s.Names) == 1 {
				continue
			}
			if decl.Doc != nil && len(decl.Specs) == 1 && len(s.Names) == 1 {
				continue
			}
			ee = append(ee, newErrNoDoc(fset, n, kindConst))
		}
	}
	return errors.Join(ee...)
}

func checkGenDeclVarValueSpec(fset *token.FileSet, decl *ast.GenDecl) error {
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
	kindConst // const
	kindFunc  // func
	kindType  // type
	kindVar   // var
	// keep-sorted end
)

func newErrNoDoc(fset *token.FileSet, id *ast.Ident, k kind) error {
	pos := fset.Position(id.Pos())
	return fmt.Errorf("%s: %s %s: %w", pos, k, id.Name, ErrNoDoc)
}

// CheckDir parses a Go package at path and ensures it has documentation set.
func CheckDir(path string) error {
	fset := token.NewFileSet()
	opts := parser.PackageClauseOnly | parser.ParseComments
	pkgs, err := parser.ParseDir(fset, path, noTest, opts)
	if err != nil {
		return err
	}
	if err := checkPackages(fset, pkgs); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func noTest(info fs.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func checkPackages(fset *token.FileSet, pkgs map[string]*ast.Package) error {
	// Must be one package
	switch len(pkgs) {
	case 0:
		return ErrNotPackage
	case 1:
		// ok
	default:
		names := strings.Join(slices.Collect(maps.Keys(pkgs)), ",")
		return fmt.Errorf("%w: %s", ErrMultiPackage, names)
	}
	for _, p := range pkgs {
		// one package
		return checkPackage(fset, p)
	}
	return nil
}

func checkPackage(fset *token.FileSet, pkg *ast.Package) error {
	var docs []*ast.CommentGroup
	for _, f := range pkg.Files {
		if f.Doc != nil {
			docs = append(docs, f.Doc)
		}
	}
	switch len(docs) {
	case 0:
		return fmt.Errorf("package %s: %w", pkg.Name, ErrNoDoc)
	case 1:
		return nil
	default:
		locs := make([]string, 0, len(docs))
		for _, d := range docs {
			locs = append(locs, fset.Position(d.Pos()).String())
		}
		sep := "\n  "
		loc := strings.Join(locs, sep)
		return fmt.Errorf("package %s: %w%s%s", pkg.Name, ErrMultiDoc, sep, loc)
	}
}
