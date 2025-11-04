// Copyright 2025 Samvel Khalatyan. All rights reserved.

package parenthesis

func Validate(s string) bool {
	stk := new(stack)
	for _, ch := range s {
		switch {
		case ch == '(' || ch == '[' || ch == '{':
			stk.push(ch)
		case ch == ')':
			if chlast, ok := stk.pop(); !ok {
				return false
			} else if chlast != '(' {
				return false
			}
		case ch == ']':
			if chlast, ok := stk.pop(); !ok {
				return false
			} else if chlast != '[' {
				return false
			}
		case ch == '}':
			if chlast, ok := stk.pop(); !ok {
				return false
			} else if chlast != '{' {
				return false
			}
		}
	}
	if !stk.empty() {
		return false
	}
	return true
}

type stack struct {
	bb []rune
}

func (s *stack) push(b rune) {
	s.bb = append(s.bb, b)
}

func (s *stack) pop() (b rune, ok bool) {
	if s.empty() {
		return b, false
	}
	n := len(s.bb) - 1
	b = s.bb[n]
	s.bb = s.bb[:n]
	return b, true
}

func (s *stack) empty() bool {
	return len(s.bb) == 0
}
