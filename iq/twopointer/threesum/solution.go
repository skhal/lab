// Copyright 2025 Samvel Khalatyan. All rights reserved.

package threesum

import (
	"maps"
	"math"
	"sort"
)

type Triplet struct {
	A, B, C int
}

func Find(nn []int) []Triplet {
	numbers := nn[:]
	sort.Ints(numbers)
	triplets := make(map[Triplet]struct{})
	for i1 := range numbers {
		x1 := numbers[i1]
		i2 := i1 + 1
		i3 := len(numbers) - 1
		for i2 < i3 {
			x2 := numbers[i2]
			x3 := numbers[i3]
			switch sum23 := x2 + x3; {
			case sum23 < -x1:
				i2 += 1
			case sum23 > -x1:
				i3 -= 1
			default:
				triplet := Triplet{x1, x2, x3}
				triplets[triplet] = struct{}{}
				i2 += 1 // search for more
			}
		}
	}
	return collectKeys(triplets)
}

func FindWithOptimizations(nn []int) []Triplet {
	if len(nn) < 3 {
		return nil
	}
	numbers := nn[:]
	sort.Ints(numbers)
	// Optimization 1: the opposite ends of the array must have opposite signs
	if math.Signbit(float64(nn[0])) == math.Signbit(float64(nn[len(nn)-1])) {
		return nil
	}
	triplets := make(map[Triplet]struct{})
	for i1, nlen := 0, len(numbers); i1 < nlen; {
		x1 := numbers[i1]
		// Optimization 2: stop when x1 becomes positive
		if x1 >= 0 {
			break
		}
		i2 := i1 + 1
		i3 := len(numbers) - 1
		for i2 < i3 {
			x2 := numbers[i2]
			x3 := numbers[i3]
			switch sum23 := x2 + x3; {
			case sum23 < -x1:
				// Optimization 4: skip the same x2 values
				for i2 += 1; i2 < i3 && x2 == numbers[i2]; i2 += 1 {
				}
			case sum23 > -x1:
				// Optimization 5: skip the same x3 values
				for i3 -= 1; i2 < i3 && x3 == numbers[i3]; i3 -= 1 {
				}
			default:
				triplet := Triplet{x1, x2, x3}
				triplets[triplet] = struct{}{}
				// Optimization 2: skip the same x2 values
				for i2 += 1; i2 < i3 && x2 == numbers[i2]; i2 += 1 {
				}
			}
		}
		// Optimization 3: Skip the same x1 values
		for i1 < nlen && numbers[i1] == x1 {
			i1 += 1
		}
	}
	return collectKeys(triplets)
}

func collectKeys(m map[Triplet]struct{}) []Triplet {
	if len(m) == 0 {
		return nil
	}
	keys := make([]Triplet, 0, len(m))
	for k := range maps.Keys(m) {
		keys = append(keys, k)
	}
	return keys
}
