// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "math/rand/v2"

type generator struct {
	num       int
	generated int

	max  Address
	addr Address
}

func newGenerator(n int, max Address) *generator {
	return &generator{
		num: n,
		max: max,
	}
}

// Next generates a new address. It returns false if the generator is exhausted,
// i.e., it has generated the number of requested addresses.
func (g *generator) Next() bool {
	if g.generated == g.num {
		return false
	}
	g.generated++
	g.addr = g.generate()
	return true
}

func (g *generator) generate() Address {
	return Address(rand.IntN(int(g.max)))
}

// Address gives access to the last generated address.
func (g *generator) Address() Address {
	return g.addr
}
