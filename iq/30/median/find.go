// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package median

import "math"

func Find(nn []int, mm []int) int {
	if len(nn) == 0 {
		return find(mm)
	}
	if len(mm) == 0 {
		return find(nn)
	}
	// pick up the shortest array
	if len(mm) < len(nn) {
		nn, mm = mm, nn
	}
	ns, ms := findSplit(nn, mm)
	a := max(ns.prev(), ms.prev())
	b := min(ns.next(), ms.next())
	if (len(nn)+len(mm))%2 == 1 {
		return b
	}
	return a + (b-a)/2
}

func find(nn []int) int {
	i := len(nn) / 2
	if len(nn)%2 == 0 {
		return nn[i-1] + (nn[i]-nn[i-1])/2
	}
	return nn[i]
}

func findSplit(nn []int, mm []int) (*split, *split) {
	halfLength := len(nn) + (len(mm)-len(nn))/2
	left := 0
	right := len(nn)
	for {
		ns := newSplit(nn, left+(right-left)/2)
		ms := newSplit(mm, halfLength-ns.idx)
		switch {
		case canMoveLeft(ns, ms):
			right = ns.idx
		case canMoveRight(ns, ms):
			left = ns.idx + 1
		default:
			return ns, ms
		}
	}
}

type split struct {
	nn  []int
	idx int
}

func newSplit(nn []int, idx int) *split {
	return &split{nn, idx}
}

func (s *split) prev() int {
	if s.idx == 0 {
		return math.MinInt
	}
	return s.nn[s.idx-1]
}

func (s *split) next() int {
	if s.idx == len(s.nn) {
		return math.MaxInt
	}
	return s.nn[s.idx]
}

func canMoveLeft(n *split, m *split) bool {
	return n.prev() > m.next()
}

func canMoveRight(n *split, m *split) bool {
	return m.prev() > n.next()
}
