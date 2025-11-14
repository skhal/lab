// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Topm finds top N transactions in the standard input.
*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"iter"
	"os"
	"slices"

	"github.com/skhal/lab/book/algos/c1/s2/fin"
	"github.com/skhal/lab/book/algos/c2/s4/queue"
)

var count = flag.Int("n", 3, "number of top transactions")

func main() {
	flag.Parse()
	if err := run(*count); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(count int) error {
	// Use Min-PQ to keep track of largest items and remove smallest of (count+1).
	pq := queue.NewUnorderedArrayPriorityQueue[*fin.Transaction](func(a, b *fin.Transaction) bool {
		return b.Amount < a.Amount
	})
	s := newTxScanner(os.Stdin)
	for tx := range s.Scan() {
		pq.Push(tx)
		if pq.Size() > count {
			pq.Pop()
		}
	}
	if err := s.Err(); err != nil {
		return err
	}
	for _, tx := range collect(pq) {
		fmt.Println(tx)
	}
	return nil
}

type TxScanner struct {
	r   io.Reader
	err error
}

func newTxScanner(r io.Reader) *TxScanner {
	return &TxScanner{
		r: r,
	}
}

func (ts *TxScanner) Scan() iter.Seq[*fin.Transaction] {
	return func(yield func(*fin.Transaction) bool) {
		for s := bufio.NewScanner(ts.r); s.Scan(); {
			tx := new(fin.Transaction)
			_, err := fmt.Sscanln(s.Text(), tx)
			if errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			if err != nil {
				ts.err = err
				break
			}
			if !yield(tx) {
				return
			}
		}
	}
}

func (ts *TxScanner) Err() error {
	return ts.err
}

type Stack interface {
	Empty() bool
	Pop()
	Top() *fin.Transaction
}

func collect(st Stack) []*fin.Transaction {
	s := slices.Collect(drain(st))
	slices.Reverse(s)
	return s
}

func drain(st Stack) iter.Seq[*fin.Transaction] {
	return func(yield func(*fin.Transaction) bool) {
		for !st.Empty() {
			tx := st.Top()
			st.Pop()
			if !yield(tx) {
				break
			}
		}
	}
}
