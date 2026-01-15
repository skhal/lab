// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Multiway merges multiples sorted files.

The files contain sorted arrays of strings:

	    % more a.txt
			a c d
			% more b.txt
			a b e
			% multiway a.txt b.txt
			a a b c d e
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/skhal/lab/book/algos/c2/s4/queue"
)

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string) error {
	pq := queue.NewMapBinaryHeapPQ[string, *bufio.Scanner](func(x, y string) bool {
		return strings.Compare(y, x) < 0
	})
	for _, fn := range filenames {
		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		s.Split(bufio.ScanWords)
		if !s.Scan() {
			continue
		}
		pq.Push(s.Text(), s)
	}
	for !pq.Empty() {
		tok, s := pq.Top()
		pq.Pop()
		fmt.Print(tok, " ")
		if !s.Scan() {
			continue
		}
		pq.Push(s.Text(), s)
	}
	fmt.Println()
	return nil
}
