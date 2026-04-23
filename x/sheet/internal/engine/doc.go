// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package engine holds different kinds of engines to drive the sheets.
//
// The engines differ in how they store parsed cell data, i.e. intermediate
// representation (IR), and how to calculate the value from the IR.
//
// For example, an AST engine may use an AST for IR. A VM engine may use an
// instruction set in the form of bytecode and calculate the result using
// virtual machines.
//
// One of the noticeable difference between the IRs is that some can be saved
// while others need to be constructed from scratch.
package engine
