// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package threesum

import (
	"fmt"
	"iter"
	"maps"
	"math"
	"sort"
)

type Triplet [3]int

func (t *Triplet) String() string {
	return fmt.Sprintf("%v", t[:])
}

func Find(nn []int) []*Triplet {
	numbers := append([]int(nil), nn...)
	sort.Ints(numbers)
	triplets := make(map[Triplet]struct{})
	for i, x := range numbers {
		for triplet := range findTwoSum(numbers[i+1:], -x) {
			triplets[*triplet] = struct{}{}
		}
	}
	return collectKeys(triplets)
}

func findTwoSum(nn []int, x int) iter.Seq[*Triplet] {
	return func(yield func(*Triplet) bool) {
		if len(nn) < 2 {
			return
		}
		i1 := 0
		i2 := len(nn) - 1
		for i1 < i2 {
			x1 := nn[i1]
			x2 := nn[i2]
			switch sum := x1 + x2; {
			case sum < x:
				i1 += 1
			case sum > x:
				i2 -= 1
			default:
				triplet := &Triplet{-x, x1, x2}
				if !yield(triplet) {
					return
				}
				i1 += 1 // search for more
			}
		}
	}
}

func FindWithOptimizations(nn []int) []*Triplet {
	if len(nn) < 3 {
		return nil
	}
	numbers := nn[:]
	sort.Ints(numbers)
	// Optimization 1: the opposite ends of the array must have opposite signs
	if !hasOppositeSignEnds(numbers) {
		return nil
	}
	triplets := make(map[Triplet]struct{})
	for i, nlen := 0, len(numbers); i < nlen; {
		x := numbers[i]
		// Optimization 2: stop when x becomes positive
		if x >= 0 {
			break
		}
		for triplet := range findTwoSumWithOptimizations(numbers[i+1:], -x) {
			triplets[*triplet] = struct{}{}
		}
		// Optimization 3: Skip the same x values
		for i < nlen && numbers[i] == x {
			i += 1
		}
	}
	return collectKeys(triplets)
}

func hasOppositeSignEnds(nn []int) bool {
	firstSign := math.Signbit(float64(nn[0]))
	lastSign := math.Signbit(float64(nn[len(nn)-1]))
	return firstSign != lastSign
}

func findTwoSumWithOptimizations(nn []int, x int) iter.Seq[*Triplet] {
	return func(yield func(*Triplet) bool) {
		if len(nn) < 2 {
			return
		}
		i1 := 0
		i2 := len(nn) - 1
		for i1 < i2 {
			x1 := nn[i1]
			x2 := nn[i2]
			switch sum := x1 + x2; {
			case sum < x:
				// Optimization 4: skip the same x1 values
				i1 += 1
				for i1 < i2 && x1 == nn[i1] {
					i1 += 1
				}
			case sum > x:
				// Optimization 5: skip the same x2 values
				i2 -= 1
				for i1 < i2 && x2 == nn[i2] {
					i2 -= 1
				}
			default:
				triplet := &Triplet{-x, x1, x2}
				if !yield(triplet) {
					return
				}
				// Optimization 4: skip the same x1 values
				i1 += 1
				for i1 < i2 && x1 == nn[i1] {
					i1 += 1
				}
			}
		}
	}
}

func collectKeys(m map[Triplet]struct{}) []*Triplet {
	if len(m) == 0 {
		return nil
	}
	keys := make([]*Triplet, 0, len(m))
	for k := range maps.Keys(m) {
		keys = append(keys, &k)
	}
	return keys
}
