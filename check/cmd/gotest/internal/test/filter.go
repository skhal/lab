// Copyright 2025 Samvel Khalatyan. All rights reserved.

package test

import (
	"iter"
	"path/filepath"
)

// FilterFunc filters items in the sequence for which fn(item) returns true.
func FilterFunc(seq iter.Seq[string], fn func(string) bool) iter.Seq[string] {
	return func(yield func(string) bool) {
		for v := range seq {
			if !fn(v) {
				continue
			}
			if !yield(v) {
				break
			}
		}
	}
}

// Pathsextracts paths from files.
func Paths(seq iter.Seq[string]) iter.Seq[string] {
	return func(yield func(string) bool) {
		for v := range seq {
			pkg, _ := filepath.Split(v)
			if !yield(pkg) {
				break
			}
		}
	}
}

// Unique suppresses strings that are equal to already passed ones.
func Unique(seq iter.Seq[string]) iter.Seq[string] {
	seen := make(map[string]struct{})
	fn := func(s string) bool {
		if _, ok := seen[s]; ok {
			return false
		}
		seen[s] = struct{}{}
		return true
	}
	return FilterFunc(seq, fn)
}
