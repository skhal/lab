// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package check verifies that Protobuf declarations are documented.
package check

import (
	"errors"
	"fmt"
	"io"

	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
)

// CheckFile validates that every item in the Protobuf, represented by the r
// reader, is documented, i.e., includes a comment. It returns a consolidated
// error of found violations.
//
// The file parameter is for Protobuf parser to prefix errors.
func CheckFile(file string, r io.Reader) error {
	f, err := parser.Parse(file, r, reporter.NewHandler(nil))
	if err != nil {
		return err
	}
	return checkFileNode(f)
}

func checkFileNode(f *ast.FileNode) error {
	var ee []error
	for _, d := range f.Decls {
		if err := checkFileElement(f, d); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func checkFileElement(f *ast.FileNode, fe ast.FileElement) error {
	switch n := fe.(type) {
	case *ast.EnumNode:
		return checkEnumNode(f, n)
	case *ast.MessageNode:
		return checkMessageNode(f, n)
	}
	return nil
}

func checkEnumNode(f *ast.FileNode, e *ast.EnumNode) error {
	prefix := func() string { return fmt.Sprintf("enum %s", e.Name.Val) }
	return checkLeadingComments(f, e, prefix)
}

func checkMessageNode(f *ast.FileNode, m *ast.MessageNode) error {
	var ee []error
	prefix := func() string { return fmt.Sprintf("message %s", m.Name.Val) }
	if err := checkLeadingComments(f, m, prefix); err != nil {
		ee = append(ee, err)
	}
	for _, c := range m.Children() {
		switch n := c.(type) {
		case *ast.EnumNode:
			if err := checkEnumNode(f, n); err != nil {
				ee = append(ee, err)
			}
		case *ast.FieldNode:
			if err := checkFieldNode(f, n); err != nil {
				ee = append(ee, err)
			}
		case *ast.MessageNode:
			if err := checkMessageNode(f, n); err != nil {
				ee = append(ee, err)
			}
		}
	}
	return errors.Join(ee...)
}

func checkFieldNode(f *ast.FileNode, n *ast.FieldNode) error {
	prefix := func() string { return fmt.Sprintf("field %s", n.Name.Val) }
	return checkLeadingComments(f, n, prefix)
}

func checkLeadingComments(f *ast.FileNode, n ast.Node, prefix func() string) error {
	if ni := f.NodeInfo(n); ni.LeadingComments().Len() == 0 {
		return fmt.Errorf("%s %s: missing comment", ni.Start(), prefix())
	}
	return nil
}
