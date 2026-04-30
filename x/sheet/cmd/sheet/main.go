// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sheet demonstrates an Excel-like engine to drive a table of cells.
//
// SYNOPSIS
//
//	sheet [-engine (ast|vm)]
//
// # DESCRIPTION
//
// Sheet is a table of cells. Every cell has a floating point number value or
// a formula. The formula can be a binary operation "=1+3" with grouping using
// parenthesis "=2 * (3 + 4)", a function call with references, ranges,
// expressions or even if-logic "=SUM(IF(A3 > 4, B2, B6), C3:C5)".
//
// Sheets eagerly parse cell contents but only calculate formulas in
// [sheet.Sheet.Calculate].
//
// There are two engines to drive parsing and calculation:
//
//   - AST: parse formulas into Abstract Syntax Tree (AST). Calculate formulas
//     by walking AST in post-traverse order.
//
//   - VM: parse formulas into AST, then compile it into an instructions set.
//     Use Virtual Machine (VM) to run the instructions.
//
// Sheet treats all values of float64 for the sake of simplicity (it is a demo
// project, actual product should support booleans, text, etc.). As such, the
// comparison operators, e.g. equal "==" and not equal "!=", use precision-based
// comparison - it checks that relative difference between two numbers is within
// pre-compiled relative value (say, 0.1%).
//
// EXAMPLE
//
//	// create an empty spreadsheet
//	s := sheet.New()
//	// fill in data
//	s.Set("A1", "123")
//	s.Set("B2", "=IF(A1 > 100, C1, D1)")
//	...
//	// calculate values, i.e. run formulas
//	s.Calculate()
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/x/sheet/internal/sheet"
)

var (
	engine  = flag.String("eng", "ast", "engine to use: ast, vm")
	engines = map[string]sheet.Option{
		"ast": sheet.WithASTEngine(),
		"vm":  sheet.WithVMEngine(),
	}
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	defer func() {
		x := recover()
		if x == nil {
			return
		}
		e, ok := x.(error)
		if !ok {
			return
		}
		err = e
	}()
	b, err := create()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return fmt.Errorf("empty sheet")
	}
	s := sheet.New()
	must(s.Read(bytes.NewReader(b)))
	must(s.Calculate())
	s.VisitAll(func(id, cell string, n float64) bool {
		fmt.Printf("%s %q = %.1f\n", id, cell, n)
		return true
	})
	return nil
}

func create() ([]byte, error) {
	var opts []sheet.Option
	if eng, ok := engines[*engine]; !ok {
		panic(fmt.Errorf("unsupported engine"))
	} else {
		opts = append(opts, eng)
	}
	s := sheet.New(opts...)
	must(s.Set("A1", "5"))
	must(s.Set("A2", "10"))
	must(s.Set("A3", "12"))
	must(s.Set("B1", "=123"))
	must(s.Set("B2", "=321"))
	must(s.Set("C1", "=1+3"))
	must(s.Set("C2", "=1-3"))
	must(s.Set("C3", "=1-3+5"))
	must(s.Set("C4", "=2*3"))
	must(s.Set("C5", "=2*3+4"))
	must(s.Set("C6", "=2+3*4"))
	must(s.Set("D1", "=(1+3)"))
	must(s.Set("D2", "=1-(2+3)"))
	must(s.Set("D3", "=1-(2-3)"))
	must(s.Set("D4", "=2*(3+4)"))
	must(s.Set("D5", "=(2+3)*4"))
	must(s.Set("E1", "=A1"))
	must(s.Set("E2", "=B1"))
	must(s.Set("E3", "=C1"))
	must(s.Set("E4", "=D2"))
	must(s.Set("E5", "=A1 - 3"))
	must(s.Set("E6", "=D2 + 10"))
	must(s.Set("F1", "=SUM()"))
	must(s.Set("F2", "=SUM(1)"))
	must(s.Set("F3", "=SUM(1, 2, 3)"))
	must(s.Set("F3", "=SUM(A1, 2)"))
	must(s.Set("F3", "=SUM(A1, A1)"))
	must(s.Set("F4", "=SUM(A1, SUM(A2, A3))"))
	must(s.Set("G1", "=SUM(A1:A3)"))
	must(s.Set("G2", "=SUM(A1:A3, 5-7)"))
	must(s.Set("G3", "=SUM(A1:A5, 1+(9-7+(2+3)), 2+4*(6-2)/(3-1)-5, B1:B5, C1:C5, D1:D5, E1:E5, 1+(2-3))"))
	must(s.Set("H1", "=IF(A1 < 10, 1, 2)"))
	must(s.Set("H2", "=IF(A1 + 5 < 10 +5, 1, 2)"))
	must(s.Set("H3", "=IF(1 < 2, 1, 1 + 2)"))
	must(s.Set("H4", "=IF(A1 + 5 < 10 +5, 1 * (3+4), 2 +3*2)"))
	must(s.Set("H5", "=IF(A1 + 5 < 10 +5, B1, B2)"))
	must(s.Set("H6", "=2*IF(A1 + 5 < 10 +5, B1, B2)"))

	var b bytes.Buffer
	must(s.Write(&b))
	return b.Bytes(), nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
