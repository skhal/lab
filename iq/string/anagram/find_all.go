// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// The solution works with Unicode code points instead of ASCII characters. This
// fact leads to increased complexity - we can't use a constant size vector to
// keep track of ASCII letter frequency as Unicode code points span larger space
// and we neeed to fall back to the hash map to track frequences.

package anagram

import (
	"iter"
	"unicode/utf8"
)

func FindAll(s string, t string) []string {
	var aa []string
	for a := range findAnagrams(s, t) {
		aa = append(aa, a)
	}
	return aa
}

func findAnagrams(s string, t string) iter.Seq[string] {
	return func(yield func(string) bool) {
		if len(s) == 0 {
			return
		}
		if len(s) < len(t) {
			return
		}
		tfp := newFootprint([]byte(t))
		for buf, fp := range substrings(s, len(t)) {
			if !tfp.Equal(fp) {
				continue
			}
			if !yield(string(buf)) {
				break
			}
		}
	}
}

type footprint map[rune]int

func newFootprint(buf []byte) footprint {
	fp := make(map[rune]int)
	for len(buf) > 0 {
		r, rs := utf8.DecodeRune(buf)
		if r == utf8.RuneError {
			break
		}
		fp[r] += 1
		buf = buf[rs:]
	}
	return footprint(fp)
}

func (fp footprint) Equal(other footprint) bool {
	if len(fp) != len(other) {
		return false
	}
	for r, n := range fp {
		nother, ok := other[r]
		if !ok {
			return false
		}
		if n != nother {
			return false
		}
	}
	return true
}

func (fp footprint) RemoveFirstRuneFrom(buf []byte) (size int, ok bool) {
	r, rs := utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		return
	}
	fp.remove(r)
	return rs, true
}

func (fp footprint) remove(r rune) {
	n, ok := fp[r]
	if !ok {
		return
	}
	n -= 1
	if n == 0 {
		delete(fp, r)
		return
	}
	fp[r] = n
}

func (fp footprint) AddLastRuneFrom(buf []byte) bool {
	r, _ := utf8.DecodeLastRune(buf)
	if r == utf8.RuneError {
		return false
	}
	fp[r] += 1
	return true
}

func substrings(s string, size int) iter.Seq2[[]byte, footprint] {
	return func(yield func([]byte, footprint) bool) {
		var fp footprint
		for buf := []byte(s); len(buf) >= size; {
			str := buf[:size]
			if fp == nil {
				fp = newFootprint(str)
			} else {
				ok := fp.AddLastRuneFrom(str)
				if !ok {
					break
				}
			}
			if !yield(str, fp) {
				break
			}
			rs, ok := fp.RemoveFirstRuneFrom(str)
			if !ok {
				break
			}
			buf = buf[rs:]
		}
	}
}
