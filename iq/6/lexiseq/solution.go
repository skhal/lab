// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lexiseq

const (
	InvalidPivot = byte(0)
	InvalidIndex = -1
)

func Next(s string) string {
	if len(s) < 2 {
		return s
	}
	buf := []byte(s[:])
	pivot, idx := findPivot(buf)
	if idx == InvalidIndex {
		reverse(buf)
		return string(buf)
	}
	nextidx := idx + 1 + findPivotNext(buf[idx+1:], pivot)
	swap(buf, idx, nextidx)
	reverse(buf[idx+1:])
	return string(buf)
}

func findPivot(bb []byte) (byte, int) {
	for j := len(bb) - 2; j >= 0; j -= 1 {
		if pivot := bb[j]; pivot < bb[j+1] {
			return pivot, j
		}
	}
	return InvalidPivot, InvalidIndex
}

func findPivotNext(bb []byte, b byte) int {
	for i := len(bb) - 1; i >= 0; i -= 1 {
		if bb[i] > b {
			return i
		}
	}
	return 0
}

func swap(bb []byte, i, j int) {
	if i == j {
		return
	}
	bb[i], bb[j] = bb[j], bb[i]
}

func reverse(bb []byte) {
	i, j := 0, len(bb)-1
	for i < j {
		swap(bb, i, j)
		i += 1
		j -= 1
	}
}
