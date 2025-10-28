// Copyright 2025 Samvel Khalatyan. All rights reserved.

package substring

import (
	"iter"
	"unicode/utf8"
)

func Find(s string, n int) string {
	var longest []byte
	for buf := range substrings(s, n) {
		if len(buf) > len(longest) {
			longest = buf
		}
	}
	return string(longest)
}

func substrings(s string, n int) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		w := newWindow(s, n)
		for !w.end() {
			w.expand()
			if !yield(w.bytes()) {
				break
			}
			w.shrink()
		}
	}
}

type window struct {
	buf   []byte
	nrepl int

	freq    map[rune]int
	freqMax int

	left  int
	right int
	count int
}

func newWindow(s string, n int) *window {
	return &window{
		buf:   []byte(s),
		nrepl: n,
		freq:  make(map[rune]int),
	}
}

func (w *window) end() bool {
	return w.right == len(w.buf)
}

func (w *window) bytes() []byte {
	return w.buf[w.left:w.right]
}

func (w *window) expand() {
	for w.right < len(w.buf) {
		r, rs := utf8.DecodeRune(w.buf[w.right:])
		if r == utf8.RuneError {
			if rs == 1 {
				// encoding error
				break
			}
			// empty
			break
		}
		if ok := w.expandOne(r, rs); !ok {
			break
		}
	}
}

func (w *window) expandOne(r rune, rs int) bool {
	canExpand := func(r rune) bool {
		rfreq := w.freq[r] + 1
		count := w.count + 1
		nmax := w.freqMax
		if rfreq > nmax {
			nmax = rfreq
		}
		replacements := count - nmax
		return replacements <= w.nrepl
	}
	if !canExpand(r) {
		return false
	}
	w.freq[r] += 1
	w.right += rs
	w.count += 1
	if n := w.freq[r]; n > w.freqMax {
		w.freqMax = n
	}
	return true
}

func (w *window) shrink() {
	r, rs := utf8.DecodeRune(w.buf[w.left:])
	if r == utf8.RuneError {
		if rs == 1 {
			// encoding error
			return
		}
		// empty
		return
	}
	w.freq[r] -= 1
	w.left += rs
	w.count -= 1
}
