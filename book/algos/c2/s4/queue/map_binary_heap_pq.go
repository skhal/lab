// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

// emptyHeap has a single, unused element at 0-index to simplify index math in
// promote() and demote().
const emptyHeapSize = 1

type keyValue[K any, V any] struct {
	key K
	val V
}

func newKeyValue[K any, V any](k K, v V) keyValue[K, V] {
	return keyValue[K, V]{k, v}
}

// MapBinaryHeapPQ is a key-value map that stores keys in a heap-oriented
// binary-tree priority queue.
type MapBinaryHeapPQ[K comparable, V any] struct {
	lessFn LessFunc[K]
	items  []keyValue[K, V]
}

// NewMapBinaryHeapPQ creates a MapBinaryHeapPQ with a given keys comparison
// function.
func NewMapBinaryHeapPQ[K comparable, V any](fn LessFunc[K]) *MapBinaryHeapPQ[K, V] {
	return &MapBinaryHeapPQ[K, V]{
		lessFn: fn,
		items:  make([]keyValue[K, V], emptyHeapSize, 2),
	}
}

// Empty reports whether the map is empty.
func (pq *MapBinaryHeapPQ[K, V]) Empty() bool {
	return pq.Size() == 0
}

// Pop removes the key-value pair with "max" key.
func (pq *MapBinaryHeapPQ[K, V]) Pop() {
	if pq.Empty() {
		return
	}
	last := len(pq.items) - 1
	pq.swap(1, last)
	pq.items = pq.items[:last]
	if pq.Empty() {
		return
	}
	pq.demote()
}

// Push inserts a new key-value pair into the map.
func (pq *MapBinaryHeapPQ[K, V]) Push(k K, v V) {
	pq.items = append(pq.items, newKeyValue(k, v))
	pq.promote()
}

// Size reports the number of key-value pairs in the map.
func (pq *MapBinaryHeapPQ[K, V]) Size() int {
	return len(pq.items) - emptyHeapSize
}

// Top returns the key-value pair with "max" key.
func (pq *MapBinaryHeapPQ[K, V]) Top() (K, V) {
	if pq.Empty() {
		var kv keyValue[K, V]
		return kv.key, kv.val
	}
	kv := pq.items[1]
	return kv.key, kv.val
}

func (pq *MapBinaryHeapPQ[K, V]) demote() {
	i := 2
	for i < len(pq.items) {
		if i+1 < len(pq.items) && pq.less(i, i+1) {
			i += 1
		}
		if !pq.less(i/2, i) {
			break
		}
		pq.swap(i/2, i)
		i *= 2
	}
}

func (pq *MapBinaryHeapPQ[K, V]) less(i, j int) bool {
	return pq.lessFn(pq.items[i].key, pq.items[j].key)
}

func (pq *MapBinaryHeapPQ[K, V]) promote() {
	i := len(pq.items) - 1
	for i > 1 && pq.less(i/2, i) {
		pq.swap(i/2, i)
		i /= 2
	}
}

func (pq *MapBinaryHeapPQ[K, V]) swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}
