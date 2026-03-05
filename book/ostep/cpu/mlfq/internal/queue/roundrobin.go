// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

import "errors"

// ErrEmpty means the queue is empty.
var ErrEmpty = errors.New("empty queue")

// RoundRobin is a queue that round robins through items. The zero value of
// RoundRobin is an empty queue ready to use.
//
// A round robin queue maintains internal iterator to the next item to
// retrieve. The iterator loops over the end and continues from the first item.
//
//		rr := new(RoundRobin)
//	 	rr.Append(1) 	// [1]
//	 	rr.Append(2) 	// [1, 2]
//		rr.Next()  		// 1
//		rr.Next()  		// 2
//		rr.Next()  		// 1
//
// Append always adds a new item to the end of the queue, regardless of the
// position to the next item:
//
//		rr := new(RoundRobin)
//	 	rr.Append(1)	// [1]
//	 	rr.Append(2)	// [1, 2]
//		rr.Next()  		// 1
//		rr.Next()  		// 2
//	 	rr.Append(3) 	// [1, 2, 3]
//		rr.Next()			// 3
//
// Pop removes the last retrieved item and shifts the next item iterator to
// the proevious item. Continuous Pop removes items in the order of round
// robin:
//
//		rr := new(RoundRobin)
//	 	rr.Append(1) // [1]
//	 	rr.Append(2) // [1, 2]
//	 	rr.Append(3) // [1, 2, 3]
//	 	rr.Next() 	 // 1
//	 	rr.Next() 	 // 2
//		rr.Pop()		 // 2
//		rr.Pop()		 // 1
//
// If the iterator reaches the beginning of the queue, Pop starts to remove
// elements in FIFO order:
//
//		rr := new(RoundRobin)
//	 	rr.Append(1) // [1]
//	 	rr.Append(2) // [1, 2]
//	 	rr.Append(3) // [1, 2, 3]
//	 	rr.Append(4) // [1, 2, 3, 4]
//	 	rr.Next() 	 // 1
//	 	rr.Next() 	 // 2
//		rr.Pop()		 // 2
//		rr.Pop()		 // 1
//		rr.Pop()		 // 3
//		rr.Pop()		 // 4
type RoundRobin struct {
	items []any
	next  int // index of the next item
}

// Append adds an item to the end of the queue.
func (q *RoundRobin) Append(v any) {
	q.items = append(q.items, v)
}

// Len returns the length of the queue.
func (q *RoundRobin) Len() int {
	return len(q.items)
}

// Next retrieves next item from the queue, starting with the oldest item.
// It panics with [ErrEmpty] if the queue is empty.
func (q *RoundRobin) Next() any {
	if len(q.items) == 0 {
		panic(ErrEmpty)
	}
	v := q.items[q.next%len(q.items)]
	q.next++
	return v
}

// Pop removes last retrieved element from the queue and decreases iterator to
// the next item. It removes items in FIFO if no items where accessed or
// iterator reaches the beginning of the queue. It panics with [ErrEmpty] if
// the queue is empty.
func (q *RoundRobin) Pop() any {
	if q.next == 0 {
		if len(q.items) == 0 {
			panic(ErrEmpty)
		}
		q.next = 1
	}
	idx := (q.next - 1) % len(q.items)
	v := q.items[idx]
	// avoid copy for the last element
	if sz := len(q.items); idx == sz-1 {
		q.items = q.items[:idx]
	} else {
		// optimize the code for production: avoid copy on every pop, instead let
		// the slice grow and copy when some threshold (allocated vs actual size)
		// is passed.
		copy(q.items[idx:], q.items[idx+1:])
		q.items = q.items[:sz-1]
	}
	q.next--
	return v
}
