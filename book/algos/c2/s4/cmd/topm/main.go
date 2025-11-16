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
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/skhal/lab/book/algos/c1/s2/fin"
	"github.com/skhal/lab/book/algos/c2/s4/queue"
)

var ErrKind = errors.New("invalid priority queue kind")

var (
	flagCount = 3
	flagKind  = KindUnorderedArray
)

func init() {
	flag.IntVar(&flagCount, "n", 3, "number of top transactions")
	vals := slices.Collect(maps.Values(kindNames))
	usage := fmt.Sprintf("priority queue kind [%s]", strings.Join(vals, ", "))
	flag.Var(&flagKind, "k", usage)
}

type Kind int

const (
	KindUnspecified Kind = iota
	// keep-sorted start
	KindBinaryHeap
	KindOrderedArray
	KindUnorderedArray
	// keep-sorted end
)

var kindNames = map[Kind]string{
	// keep-sorted start
	KindBinaryHeap:     "binary-heap",
	KindOrderedArray:   "ordered-array",
	KindUnorderedArray: "unordered-array",
	KindUnspecified:    "unspecified",
	// keep-sorted end
}

func (k Kind) String() string {
	name, ok := kindNames[k]
	if !ok {
		return "unknown"
	}
	return name
}

func (k *Kind) Set(s string) error {
	switch s {
	default:
		return ErrKind
	case "unordered-array":
		*k = KindUnorderedArray
	case "ordered-array":
		*k = KindOrderedArray
	case "binary-heap":
		*k = KindBinaryHeap
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(flagCount, flagKind); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type PriorityQueue interface {
	// keep-sorted start
	Empty() bool
	Pop()
	Push(*fin.Transaction)
	Size() int
	Top() *fin.Transaction
	// keep-sorted end
}

type LessFunc queue.LessFunc[*fin.Transaction]
type NewPQFunc func(LessFunc) PriorityQueue

var newPQFuncs = map[Kind]NewPQFunc{
	KindUnorderedArray: func(less LessFunc) PriorityQueue {
		l := queue.LessFunc[*fin.Transaction](less)
		return queue.NewUnorderedArrayPQ(l)
	},
	KindOrderedArray: func(less LessFunc) PriorityQueue {
		l := queue.LessFunc[*fin.Transaction](less)
		return queue.NewOrderedArrayPQ(l)
	},
	KindBinaryHeap: func(less LessFunc) PriorityQueue {
		l := queue.LessFunc[*fin.Transaction](less)
		return queue.NewBinaryHeapPQ(l)
	},
}

// Use Min-PQ to keep track of largest items and remove smallest.
func less(a, b *fin.Transaction) bool {
	return b.Amount < a.Amount
}

func run(count int, kind Kind) error {
	f, ok := newPQFuncs[kind]
	if !ok {
		return ErrKind
	}
	pq := f(less)
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
