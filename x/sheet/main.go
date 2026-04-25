// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sheet demonstrates a cells table engine.
//
// SYNOPSIS
//
//	sheet
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
		fmt.Printf("%s %q = %.2f\n", id, cell, n)
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

	var b bytes.Buffer
	must(s.Write(&b))
	return b.Bytes(), nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
