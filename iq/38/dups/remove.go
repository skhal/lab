// Copyright 2025 Samvel Khalatyan. All rights reserved.

package dups

import "unicode/utf8"

func Remove(s string) string {
	var rr stack
	for i, buf := 0, []byte(s); i < len(buf); {
		r, rs := utf8.DecodeRune(buf[i:])
		if r == utf8.RuneError {
			break
		}
		i += rs
		if r == rr.Top() {
			rr = rr.Pop()
			continue
		}
		rr = append(rr, r)
	}
	return string(rr)
}

type stack []rune

func (s stack) Empty() bool {
	return len(s) == 0
}

func (s stack) Top() rune {
	if s.Empty() {
		return utf8.RuneError
	}
	return s[len(s)-1]
}

func (s stack) Pop() stack {
	return s[:len(s)-1]
}
