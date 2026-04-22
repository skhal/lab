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
	"runtime/pprof"

	"github.com/skhal/lab/x/sheet/internal/sheet"
)

var cpuprofile = flag.String("cpuprofile", "", "write CPU profile to file")

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
	if *cpuprofile != "" {
		var f *os.File
		f, err = os.Create(*cpuprofile)
		if err != nil {
			return fmt.Errorf("failed to create CPU profile: %s", err)
		}
		defer f.Close()
		if err = pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("failed to start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
	}
	b, err := create()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return fmt.Errorf("empty sheet")
	}
	s := sheet.New()
	err = s.Read(bytes.NewReader(b))
	if err != nil {
		return err
	}
	must(s.Calculate())
	s.VisitAll(func(id, cell string, n float64) bool {
		fmt.Printf("%s %q = %.2f\n", id, cell, n)
		return true
	})
	return nil
}

func create() ([]byte, error) {
	s := sheet.New()
	must(s.Set("A1", "5"))
	must(s.Set("A2", "10"))
	must(s.Set("A3", "12"))
	must(s.Set("B1", "=123"))
	must(s.Set("C1", "=1+3"))
	must(s.Set("C2", "=1-3"))
	must(s.Set("C3", "=1-3+5"))
	must(s.Set("D1", "=(1+3)"))
	must(s.Set("D2", "=1-(2+3)"))
	must(s.Set("D3", "=1-(2-3)"))
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
	must(s.Set("G3", "=SUM(A1:A5, 1+(9-7+(2+3)), B1:B5, C1:C5, D1:D5, E1:E5, 1+(2-3))"))
	var b bytes.Buffer
	if err := s.Write(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
