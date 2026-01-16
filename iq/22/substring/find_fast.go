// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package substring

import "unicode/utf8"

func FindFast(s string) string {
	return findWith(newFastWindow(s))
}

type fastWindow struct {
	buf  []byte
	seen map[rune]int

	left  int
	right int

	stop rune
}

func newFastWindow(s string) *fastWindow {
	return &fastWindow{
		buf:  []byte(s),
		seen: make(map[rune]int),
	}
}

func (w *fastWindow) end() bool {
	return w.left == len(w.buf)
}

func (w *fastWindow) bytes() []byte {
	return w.buf[w.left:w.right]
}

func (w *fastWindow) expand() {
	for w.left < len(w.buf) {
		r, rs := utf8.DecodeRune(w.buf[w.right:])
		if r == utf8.RuneError {
			w.seen[r] = w.right
			w.stop = r
			if rs == 1 {
				// encoding error
				break
			}
			// empty buffer
			break
		}
		if i, ok := w.seen[r]; ok && w.left <= i {
			// duplicate
			w.stop = r
			break
		}
		w.seen[r] = w.right
		w.right += rs
	}
}

func (w *fastWindow) shrink() {
	w.left = w.seen[w.stop]
	r, rs := utf8.DecodeRune(w.buf[w.left:])
	if r == utf8.RuneError {
		if rs == 1 {
			// encoding error
			return
		}
		// empty buffer
		w.left = w.right
		return
	}
	w.left += rs
}
