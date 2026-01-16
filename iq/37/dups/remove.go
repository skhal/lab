// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dups

import (
	"unicode/utf8"
)

func Remove(s string) string {
	buf := []byte(s)
	var rr stack
	for idx := 0; idx < len(buf); {
		r, rs := utf8.DecodeRune(buf[idx:])
		if r == utf8.RuneError {
			break
		}
		idx += rs
		ds := skipDuplicates(buf[idx:], r)
		idx += ds
		rrlen := len(rr)
		rr = rr.popFunc(func(top rune) bool { return top == r })
		if rrlen != len(rr) || ds > 0 {
			continue
		}
		rr = append(rr, r)
	}
	return string(rr)
}

func skipDuplicates(data []byte, r rune) int {
	idx := 0
	for idx < len(data) {
		nr, nrs := utf8.DecodeRune(data[idx:])
		if nr == utf8.RuneError {
			break
		}
		if nr != r {
			break
		}
		idx += nrs
	}
	return idx
}

type stack []rune

func (s stack) empty() bool {
	return len(s) == 0
}

func (s stack) top() rune {
	return s[len(s)-1]
}

func (s stack) pop() stack {
	return s[:len(s)-1]
}

func (s stack) popFunc(f func(rune) bool) stack {
	for !s.empty() && f(s.top()) {
		s = s.pop()
	}
	return s
}
