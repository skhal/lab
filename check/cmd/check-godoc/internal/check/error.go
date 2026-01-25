// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

// ErrNoDoc represents missing documentation error.
var ErrNoDoc = errors.New("no documentation")

// ErrNotPackage means that Go failed to parse a path.
var ErrNotPackage = errors.New("not a package")

// ErrMultiPackage means there are multiple packages at a single path.
var ErrMultiPackage = errors.New("multiple packages")

// ErrMultiDoc means there are multiple package documentations.
var ErrMultiDoc = errors.New("multiple documentation")

type kind int

//go:generate stringer -type=kind -linecomment
const (
	_ kind = iota
	// keep-sorted start
	kindConst  // const
	kindField  // field
	kindFunc   // func
	kindMethod // method
	kindType   // type
	kindVar    // var
	// keep-sorted end
)

func newErrNoDoc(fset *token.FileSet, id *ast.Ident, k kind) error {
	pos := fset.Position(id.Pos())
	return fmt.Errorf("%s: %s %s: %w", pos, k, id.Name, ErrNoDoc)
}
