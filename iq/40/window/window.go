// Copyright 2025 Samvel Khalatyan. All rights reserved.

package window

import "errors"

const MinSize = 1
const NumError = 0

type W struct {
	nn   []int
	size int

	start int
	end   int

	mm dequeue
}

func New(nn []int, size int) (*W, error) {
	if size < MinSize {
		return nil, errors.New("window size must be a positive integer")
	}
	w := &W{
		nn:   nn,
		size: size,
	}
	return w, nil
}

func (w *W) Max() int {
	if w.mm.empty() {
		return NumError
	}
	return w.mm.front()
}

func (w *W) Slide() bool {
	if w.end == len(w.nn) {
		return false
	}
	w.rebalance()
	w.includeEnd()
	return true
}

func (w *W) full() bool {
	return w.end - w.start == w.size
}

func (w *W) includeEnd() {
	w.mm = append(w.mm, item{
		num: w.nn[w.end],
		pos: w.end,
	})
	w.end += 1
}

func (w *W) popStart() {
	if !w.full() {
		return
	}
	if w.start == w.mm[0].pos {
		w.mm = w.mm.popFront()
	}
	w.start += 1
}

func (w *W) rebalance() {
	if w.end == 0 {
		return
	}
	w.popStart()
	for n := w.nn[w.end]; !w.mm.empty() && w.mm.back() <= n; {
		w.mm = w.mm.popBack()
	}
}

type item struct {
	num int
	pos int
}

type dequeue []item

func (d dequeue) empty() bool {
	return len(d) == 0
}

func (d dequeue) front() int {
	return d[0].num
}

func (d dequeue) popBack() dequeue {
	return d[:len(d)-1]
}

func (d dequeue) popFront() dequeue {
	return d[1:]
}

func (d dequeue) back() int {
	return d[len(d)-1].num
}
