// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

type Q struct {
	in  stack
	out stack
}

func New(nn ...int) *Q {
	return &Q{
		in: append(stack(nil), nn...),
	}
}

func (q *Q) Empty() bool {
	return q.Size() == 0
}

func (q *Q) Front() (int, bool) {
	if q.Empty() {
		return 0, false
	}
	if q.out.empty() {
		q.moveInToOut()
	}
	return q.out.top(), true
}

func (q *Q) Pop() {
	if q.Empty() {
		return
	}
	if q.out.empty() {
		q.moveInToOut()
	}
	q.out = q.out.pop()
}

func (q *Q) Push(n int) {
	q.in = append(q.in, n)
}

func (q *Q) Size() int {
	return len(q.in) + len(q.out)
}

func (q *Q) moveInToOut() {
	for !q.in.empty() {
		q.out = append(q.out, q.in.top())
		q.in = q.in.pop()
	}
}

type stack []int

func (s stack) empty() bool {
	return len(s) == 0
}

func (s stack) pop() stack {
	return s[:len(s)-1]
}

func (s stack) top() int {
	return s[len(s)-1]
}
