// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Life emulates the Game of Life.
//
// # SYNOPSIS
//
//	life -n 100 -delay 50ms
//
// # DESCRIPTION
//
// Game of Life runs on a binary board, where every cell can be either in ON or
// OFF state, i.e. alive or dead. In every cycle, the new world generates
// according to the following rules:
//
//  1. ON -> OFF if there are less than 2 or more than 3 neighbours.
//  2. OFF -> ON if there are 3 alive neighbours
//
// where neighbours are cells next to the currently analyzed at position x, y.
// Out of bounds cells are not considered.
//
// https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"time"
)

const initWorld = `
.................
.................
....xxx...xxx....
.................
..x....x.x....x..
..x....x.x....x..
..x....x.x....x..
....xxx...xxx....
.................
....xxx...xxx....
..x....x.x....x..
..x....x.x....x..
..x....x.x....x..
.................
....xxx...xxx....
.................
.................
`

const (
	off = '.'
	on  = 'x'
	eol = '\n'
)

var (
	cycles = flag.Int("n", 10, "life cycles to run")
	delay  = flag.Duration("delay", 40*time.Millisecond, "delay between cycles")
)

func main() {
	flag.Parse()

	w := newWorld(initWorld)
	for range *cycles {
		w.Print()
		w.Generate()
		time.Sleep(*delay)
	}
}

type world struct {
	w   int
	h   int
	now [][]byte
	buf [][]byte // buffer that is used to generate the next world
}

func newWorld(s string) *world {
	var w, h, idx int
	var (
		d   [][]byte
		buf [][]byte
	)
	for ln := range bytes.Lines(bytes.TrimSpace([]byte(s))) {
		h++
		ln = bytes.TrimSpace(ln)
		switch n := len(ln); {
		case w == 0:
			w = n
		case w != n:
			panic(fmt.Errorf("world:%d: wrong line length %d, want %d", h, n, w))
		}
		d = append(d, ln)
		buf = append(buf, make([]byte, w))
		idx += w
	}
	return &world{w, h, d, buf}
}

// Print dumps the world to standard output after clearing the screen.
func (w *world) Print() {
	const (
		clearScreen = "\033[H\033[2J"
	)
	var s strings.Builder
	for y := range w.h {
		s.Write(w.now[y])
		s.WriteByte(eol)
	}
	fmt.Print(clearScreen, s.String())
}

// Generate creates a new world from the current state.
func (w *world) Generate() {
	for y := range w.h {
		for x := range w.w {
			cell := w.now[y][x]
			switch n := w.neighbors(x, y); {
			case cell != off && n < 2 || n > 3:
				cell = off
			case cell == off && n == 3:
				cell = on
			case cell != off:
				cell = on
			default:
				cell = off
			}
			w.buf[y][x] = cell
		}
	}
	w.now, w.buf = w.buf, w.now
}

func (w *world) neighbors(x, y int) int {
	n := 0
	for dy := -1; dy < 2; dy++ {
		ny := y + dy
		if ny < 0 || ny >= w.h {
			continue
		}
		for dx := -1; dx < 2; dx++ {
			nx := x + dx
			if nx < 0 || nx >= w.w {
				continue
			}
			switch {
			case nx == x && ny == y:
				// self
			case w.now[ny][nx] != off:
				n++
			}
		}
	}
	return n
}
