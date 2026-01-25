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

// CheckFile verifies that non-generated Go file has documentation attached to
// the exported declarations. It returns an error if the check fails.
func CheckFile(fname string) error {
	fs := token.NewFileSet()
	opts := parser.SkipObjectResolution | parser.ParseComments
	f, err := parser.ParseFile(fs, fname, nil, opts)
	if err != nil {
		return err
	}
	return CheckAST(fs, f)
}

// CheckAST scans top-level exported declarations for presence of documentation.
// It skips generated data.
func CheckAST(fs *token.FileSet, f *ast.File) error {
	if ast.IsGenerated(f) {
		return nil
	}
	return checkDecls(fs, f)
}

func checkDecls(fs *token.FileSet, f *ast.File) error {
	var ee []error
	for _, d := range f.Decls {
		if err := checkDecl(fs, d); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func checkDecl(fs *token.FileSet, decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return checkFuncDecl(fs, d)
	case *ast.GenDecl:
		return checkGenDecl(fs, d)
	}
	return nil
}

func checkFuncDecl(fs *token.FileSet, decl *ast.FuncDecl) error {
	if !decl.Name.IsExported() {
		return nil
	}
	if decl.Doc != nil {
		return nil
	}
	return newErrNoDoc(fs, decl.Name, kindFunc)
}

func checkGenDecl(fs *token.FileSet, decl *ast.GenDecl) error {
	switch decl.Tok {
	case token.CONST:
		return checkGenDeclConstValueSpec(fs, decl)
	case token.TYPE:
		return checkGenDeclTypeSpec(fs, decl)
	case token.VAR:
		return checkGenDeclVarValueSpec(fs, decl)
	}
	return nil
}
func checkGenDeclConstValueSpec(fs *token.FileSet, decl *ast.GenDecl) error {
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
			ee = append(ee, newErrNoDoc(fs, n, kindConst))
		}
	}
	return errors.Join(ee...)
}

func checkGenDeclVarValueSpec(fs *token.FileSet, decl *ast.GenDecl) error {
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
			ee = append(ee, newErrNoDoc(fs, n, kindVar))
		}
	}
	return errors.Join(ee...)
}

func checkGenDeclTypeSpec(fs *token.FileSet, decl *ast.GenDecl) error {
	var ee []error
	for _, spec := range decl.Specs {
		s := spec.(*ast.TypeSpec)
		if !s.Name.IsExported() {
			continue
		}
		if err := checkTypeSpec(fs, decl, s); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func checkTypeSpec(fs *token.FileSet, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	if err := checkTypeSpecDoc(fs, decl, spec); err != nil {
		return err
	}
	return checkTypeSpecFields(fs, spec)
}

func checkTypeSpecDoc(fs *token.FileSet, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	// A comment attached to a struct in a type group:
	//  type (
	//    // comment
	//    A struct {}
	//  )
	if spec.Doc != nil {
		return nil
	}
	// A comment attached to the type group, allow only one type inside:
	//  // comment
	//  type A struct {}
	if decl.Doc != nil && len(decl.Specs) == 1 {
		return nil
	}
	return newErrNoDoc(fs, spec.Name, kindType)
}

func checkTypeSpecFields(fs *token.FileSet, ts *ast.TypeSpec) error {
	switch t := ts.Type.(type) {
	case *ast.StructType:
		return checkStructType(fs, t)
	case *ast.InterfaceType:
		return checkInterfaceType(fs, t)
	}
	return nil
}

func checkStructType(fs *token.FileSet, st *ast.StructType) error {
	// Examples:
	// - field doc: https://pkg.go.dev/runtime/metrics#Float64Histogram
	// - field comment: https://pkg.go.dev/strconv#NumError
	var ee []error
	for _, f := range st.Fields.List {
		if f.Doc != nil && len(f.Names) == 1 {
			continue
		}
		if f.Comment != nil && len(f.Names) == 1 {
			continue
		}
		for _, n := range f.Names {
			if !n.IsExported() {
				continue
			}
			ee = append(ee, newErrNoDoc(fs, n, kindField))
		}
	}
	return errors.Join(ee...)
}

func checkInterfaceType(fs *token.FileSet, it *ast.InterfaceType) error {
	// Examples:
	// - doc: https://pkg.go.dev/context#Context
	// - comment: https://pkg.go.dev/go/types#Object
	var ee []error
	for _, f := range it.Methods.List {
		if f.Doc != nil && len(f.Names) == 1 {
			continue
		}
		if f.Comment != nil && len(f.Names) == 1 {
			continue
		}
		for _, n := range f.Names {
			if !n.IsExported() {
				continue
			}
			ee = append(ee, newErrNoDoc(fs, n, kindMethod))
		}
	}
	return errors.Join(ee...)
}

// CheckDir parses a Go package at path and ensures it has documentation set.
func CheckDir(path string) error {
	fs := token.NewFileSet()
	opts := parser.PackageClauseOnly | parser.ParseComments
	pkgs, err := parser.ParseDir(fs, path, noTest, opts)
	if err != nil {
		return err
	}
	if err := checkPackages(fs, pkgs); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func noTest(info fs.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func checkPackages(fs *token.FileSet, pkgs map[string]*ast.Package) error {
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
		return checkPackage(fs, p)
	}
	return nil
}

func checkPackage(fs *token.FileSet, pkg *ast.Package) error {
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
			locs = append(locs, fs.Position(d.Pos()).String())
		}
		sep := "\n  "
		loc := strings.Join(locs, sep)
		return fmt.Errorf("package %s: %w%s%s", pkg.Name, ErrMultiDoc, sep, loc)
	}
}
