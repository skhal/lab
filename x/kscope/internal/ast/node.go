// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "fmt"

// Node is any AST node.
type Node any

// Number is a number literal node.
type Number struct {
	Val float64 // parsed value of the number.
}

// String prints the number.
func (n Number) String() string { return fmt.Sprintf("%.1f", n.Val) }
