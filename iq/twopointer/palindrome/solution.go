// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Note: the solution uses unicode code points instead of bytes (aka chars).

package palindrome

import (
	"iter"
	"unicode"
	"unicode/utf8"
)

func Is(s string) bool {
	for pair := range findRunePairs(s) {
		if !isEqual(pair.left, pair.right, unicode.ToLower) {
			return false
		}
	}
	return true
}

func findRunePairs(s string) iter.Seq[*pair] {
	return func(yield func(*pair) bool) {
		bb := []byte(s)
		i1 := 0
		i2 := len(bb)
		for i1 < i2 {
			r1, offset1 := decodeRune(bb[i1:i2], isAlphanum)
			if r1 == utf8.RuneError {
				return
			}
			r2, offset2 := decodeLastRune(bb[i1+offset1:i2], isAlphanum)
			if r2 == utf8.RuneError {
				return
			}
			if !yield(&pair{r1, r2}) {
				return
			}
			i1 += offset1
			i2 -= offset2
		}
	}
}

func decodeRune(bb []byte, f func(rune) bool) (r rune, offset int) {
	var size int
	for offset < len(bb) {
		r, size = utf8.DecodeRune(bb[offset:])
		offset += size
		if f(r) {
			return
		}
	}
	return utf8.RuneError, 0
}

func decodeLastRune(bb []byte, f func(rune) bool) (r rune, offset int) {
	var size int
	for offset < len(bb) {
		r, size = utf8.DecodeLastRune(bb[:len(bb)-offset])
		offset += size
		if f(r) {
			return
		}
	}
	return utf8.RuneError, 0
}

func isAlphanum(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

func isEqual(r1, r2 rune, f func(rune) rune) bool {
	return f(r1) == f(r2)
}

type pair struct {
	left, right rune
}
