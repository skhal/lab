// Copyright 2025 Samvel Khalatyan. All rights reserved.

package calc

import (
	"errors"
	"fmt"
)

var ErrExpression = errors.New("invalid expression")

func Eval(s string) (int, error) {
	e := new(evaluator)
	return e.evaluate([]byte(s))
}

type evaluator struct {
	ops  stack[byte] // operators
	nums stack[int]  // numbers
}

func (e *evaluator) evaluate(buf []byte) (int, error) {
	for idx := 0; idx < len(buf); {
		if idx += skipWhitespace(buf[idx:]); idx == len(buf) {
			break
		}
		ch := buf[idx]
		switch {
		case isDigit(ch):
			k, ks := parseInt(buf[idx:])
			switch {
			case e.ops.empty():
				// operators must separate numbers
				if !e.nums.empty() {
					return 0, newExpressionError(buf, idx, "missing operator")
				}
			case e.ops.top() != '(':
				res, err := e.evalOperatorWithRight(k)
				if err != nil {
					return 0, newExpressionError(buf, idx, err.Error())
				}
				k = res
			}
			e.nums = e.nums.push(k)
			idx += ks
		case isOperator(ch):
			e.ops = e.ops.push(ch)
			idx += 1
		case ch == '(':
			e.ops = e.ops.push(ch)
			idx += 1
		case ch == ')':
			if err := e.evalParenthesis(); err != nil {
				return 0, newExpressionError(buf, idx, err.Error())
			}
			if err := e.evalOperators(); err != nil {
				return 0, newExpressionError(buf, idx, err.Error())
			}
			idx += 1
		default:
			return 0, newExpressionError(buf, idx, "invalid character")
		}
	}
	switch {
	case !e.ops.empty():
		return 0, ErrExpression
	case e.nums.empty():
		return 0, nil
	case len(e.nums) > 1:
		return 0, ErrExpression
	}
	return e.nums.top(), nil
}

func (e *evaluator) evalOperatorWithRight(n int) (int, error) {
	if e.nums.empty() {
		return 0, fmt.Errorf("missing left number for operator")
	}
	n, err := evalOperator(e.ops.top(), e.nums.top(), n)
	if err != nil {
		return 0, err
	}
	e.ops = e.ops.pop()
	e.nums = e.nums.pop()
	return n, nil
}

func evalOperator(op byte, x, y int) (int, error) {
	switch op {
	case '+':
		x = x + y
	case '-':
		x = x - y
	default:
		return 0, fmt.Errorf("unsupported operator %c", op)
	}
	return x, nil
}

func (e *evaluator) evalOperators() error {
	n := e.nums.top()
	e.nums = e.nums.pop()
	for !e.ops.empty() && e.ops.top() != '(' {
		res, err := e.evalOperatorWithRight(n)
		if err != nil {
			return err
		}
		n = res
	}
	e.nums = e.nums.push(n)
	return nil
}

func (e *evaluator) evalParenthesis() error {
	if e.ops.empty() {
		return errors.New("missing parenthesis")
	}
	if e.ops.top() != '(' {
		return errors.New("unexpected operator")
	}
	e.ops = e.ops.pop()
	// prevent empty parenthesis
	if e.nums.empty() {
		return errors.New("missing numbers")
	}
	return nil
}

func skipWhitespace(buf []byte) (n int) {
	for n < len(buf) && isWhitespace(buf[n]) {
		n += 1
	}
	return
}

func isWhitespace(b byte) bool {
	switch b {
	case ' ':
	case '\t':
	default:
		return false
	}
	return true
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isOperator(b byte) bool {
	switch b {
	case '+':
	case '-':
	default:
		return false
	}
	return true
}

func parseInt(buf []byte) (n int, ns int) {
	for ns < len(buf) {
		ch := buf[ns]
		if !isDigit(ch) {
			break
		}
		n = n*10 + int(ch-'0')
		ns += 1
	}
	return
}

type Item interface {
	byte | int
}

type stack[V Item] []V

func (s stack[V]) empty() bool {
	return len(s) == 0
}

func (s stack[V]) pop() stack[V] {
	return s[:len(s)-1]
}

func (s stack[V]) top() V {
	return s[len(s)-1]
}

func (s stack[V]) push(v V) stack[V] {
	return append(s, v)
}

type expressionError struct {
	buf []byte
	pos int
	msg string
}

func newExpressionError(buf []byte, pos int, msg string) *expressionError {
	return &expressionError{
		buf: buf,
		pos: pos,
		msg: msg,
	}
}

func (e *expressionError) Is(err error) bool {
	return err == ErrExpression
}

func (e *expressionError) Error() string {
	return fmt.Sprintf("%s: %q at %q: %s", ErrExpression, e.buf, e.buf[e.pos:], e.msg)
}
