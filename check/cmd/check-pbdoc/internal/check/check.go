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
	"strings"

	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
)

var (
	// ErrCommentPrefix means the comments starts with a wrong prefix.
	// Message or enumeration documentation must start with the item name:
	//	// Foo ...
	// 	message Foo {}
	ErrCommentPrefix = errors.New("wrong comment prefix")

	// ErrNoComment means no documentation found.
	ErrNoComment = errors.New("missing documentation")
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
	var ee []error
	prefix := func() string { return fmt.Sprintf("enum %s", e.Name.Val) }
	if err := checkLeadingComments(f, e, prefix); err != nil {
		ee = append(ee, err)
	}
	for _, c := range e.Children() {
		switch n := c.(type) {
		case *ast.EnumValueNode:
			if err := checkEnumValueNode(f, n); err != nil {
				ee = append(ee, err)
			}
		}
	}
	return errors.Join(ee...)
}

func checkEnumValueNode(f *ast.FileNode, n *ast.EnumValueNode) error {
	prefix := func() string { return fmt.Sprintf("enumerator %s", n.Name.Val) }
	return checkLeadingComments(f, n, prefix)
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
	ni := f.NodeInfo(n)
	if ni.LeadingComments().Len() == 0 {
		return fmt.Errorf("%s %s: %w", ni.Start(), prefix(), ErrNoComment)
	}
	switch node := n.(type) {
	case *ast.EnumNode:
		if err := checkLeadingCommentsPrefix(ni, node.Name.Val); err != nil {
			return fmt.Errorf("%s %s: %w", ni.Start(), prefix(), err)
		}
	case *ast.MessageNode:
		if err := checkLeadingCommentsPrefix(ni, node.Name.Val); err != nil {
			return fmt.Errorf("%s %s: %w", ni.Start(), prefix(), err)
		}
	}
	return nil
}

func checkLeadingCommentsPrefix(ni ast.NodeInfo, name string) error {
	p := fmt.Sprintf("// %s", name)
	if !strings.HasPrefix(ni.LeadingComments().Index(0).RawText(), p) {
		return fmt.Errorf("%w, want %q", ErrCommentPrefix, p)
	}
	return nil
}
