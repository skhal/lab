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
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/sheet"
)

func main() {
	s := sheet.New()
	s.Set("A1", "5")
	s.Set("A2", "10")
	s.Set("A3", "12")
	s.Calculate()
	s.VisitAll(func(c string, n float64) bool {
		fmt.Printf("%s: %.2f\n", c, n)
		return true
	})
}
