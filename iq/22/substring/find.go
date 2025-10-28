// Copyright 2025 Samvel Khalatyan. All rights reserved.

package substring

import (
	"iter"
	"unicode/utf8"
)

func Find(s string) string {
	var longest []byte
	for buf := range findUnique(s) {
		if len(buf) > len(longest) {
			longest = buf
		}
	}
	return string(longest)
}

type window struct {
	buf   []byte
	seen  map[rune]struct{}

	left  int
	right int
	stop  rune
}

func newWindow(s string) *window {
	return &window{
		buf:  []byte(s),
		seen: make(map[rune]struct{}),
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
			w.stop = r
			if rs == 1 {
				// encoding error
				break
			}
			// end of stream
			break
		}
		if _, ok := w.seen[r]; ok {
			w.stop = r
			break
		}
		w.seen[r] = struct{}{}
		w.right += rs
	}
}

func (w *window) shrink() {
	for w.left < w.right {
		r, rs := utf8.DecodeRune(w.buf[w.left:])
		if r == utf8.RuneError {
			if rs == 1 {
				// encoding error
				break
			}
			// end of stream
			break
		}
		delete(w.seen, r)
		w.left += rs
		if r == w.stop {
			w.stop = utf8.RuneError
			break
		}
	}
}

func findUnique(s string) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		w := newWindow(s)
		for !w.end() {
			w.expand()
			if !yield(w.bytes()) {
				break
			}
			w.shrink()
		}
	}
}
