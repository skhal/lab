// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package largest

const nextError = -1

func Find(nn []int) []int {
	mm := make([]int, len(nn))
	var s stack
	for i := len(nn); i > 0; {
		i -= 1
		n := nn[i]
		for !s.empty() {
			if s.peek() > n {
				break
			}
			s = s.pop()
		}
		mm[i] = s.top()
		s = append(s, n)
	}
	return mm
}

type stack []int

func (s stack) empty() bool {
	return len(s) == 0
}

func (s stack) top() int {
	if s.empty() {
		return nextError
	}
	return s.peek()
}

func (s stack) peek() int {
	return s[len(s)-1]
}

func (s stack) pop() stack {
	return s[:len(s)-1]
}
