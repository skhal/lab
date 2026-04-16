// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"unicode/utf8"
)

// ErrLex means lexican analysis failed.
var ErrLex = errors.New("lex error")

// Lex runs lexical analysis on the formula and generates a sequence of tokens.
func Lex(b []byte) iter.Seq[Token] {
	return newLexer(b).Run()
}

// lexer runs lexical analysis on the formula. The implementation borrows the
// idea of concurrent "Lexical Scanning" from Rob Pike's presentation:
//
//	https://go.dev/talks/2011/lex.slide
//
// At high level:
//
//   - The lexer runs in a separate Go-routine.
//
//   - Lexer state machine uses transitive states: every state knows the next
//     state it may transition into. For example, formula "1 + 2" has following
//     states:
//     [scan, number, scan, plus, scan, number, eof]
//     where scan-state skips whitespace and identifies the next state.
//
//   - Some states emit designatited tokens, e.g. a number-state emits
//     TokenNumber.
type lexer struct {
	b    []byte
	toks chan Token
	done chan struct{}

	// state
	start        int // starting position of the token being scanned
	pos          int // current positin during the scan
	lastRuneSize int
	err          error
}

func newLexer(b []byte) *lexer {
	return &lexer{
		b:    b,
		toks: make(chan Token),
		done: make(chan struct{}),
	}
}

// Run runs state machine staring with scanState. It generates a sequence of
// tokens. If lexer detects an error, it emits TokenError and stops the
// sequence.
// Note: call only once per lexer instance.
func (lx *lexer) Run() iter.Seq[Token] {
	go func() {
		defer close(lx.toks)
		for st := scanState; st != nil; st = st(lx) {
			select {
			case <-lx.done:
				return
			default:
			}
		}
	}()
	return func(yield func(Token) bool) {
		for {
			tok, ok := <-lx.toks
			if !ok {
				break
			}
			if !yield(tok) {
				lx.Stop()
				break
			}
		}
	}
}

// Stop ends the lexer. Call it once.
func (lx *lexer) Stop() {
	close(lx.done)
}

// emit creates a token and sends it on the tokkens channel. It blocks until
// token receiver reads the token or if [lexer.Stop] called.
func (lx *lexer) emit(tk tokenType) {
	var tok Token
	switch tk {
	case TokenError:
		err := fmt.Errorf("%w: %s", ErrLex, lx.err)
		tok = Token{Type: tk, Err: err}
	default:
		tok = Token{Type: tk, Text: lx.text()}
	}
	select {
	case <-lx.done:
	case lx.toks <- tok:
	}
}

// text returns token text [start:pos) and calls [lexer.ignore].
func (lx *lexer) text() string {
	b := lx.b[lx.start:lx.pos]
	lx.ignore()
	return string(b)
}

// ignore brings start to current positin and resets the size of the last read
// rune.
func (lx *lexer) ignore() {
	lx.start = lx.pos
	lx.lastRuneSize = 0
}

// peek returns the next rune in the buffer without advancing current position.
func (lx *lexer) peek() (r rune, err error) {
	r, _, err = lx.decodeRune()
	if err != nil {
		return
	}
	return
}

// decodeRune reads next rune without advancing current position. It returns
// decoded rune and its size upon success, owtherwise an error.
func (lx *lexer) decodeRune() (r rune, sz int, err error) {
	if lx.pos == len(lx.b) {
		return r, sz, io.EOF
	}
	r, sz = utf8.DecodeRune(lx.b[lx.pos:])
	if r == utf8.RuneError {
		switch sz {
		case 0:
			r = rune(0)
			err = io.EOF
			return
		case 1:
			r = rune(0)
			sz = 0
			err = fmt.Errorf("encoding error at %d - %q", lx.pos, lx.b)
			return
		}
	}
	return
}

// read reads the next rune and advances current position in the buffer. It
// also sets the size of the last read rune to unblock [lexer.unread].
func (lx *lexer) read() (r rune, err error) {
	r, sz, err := lx.decodeRune()
	if err != nil {
		return
	}
	lx.pos += sz
	lx.lastRuneSize = sz
	return r, nil
}

// unread reverses the position in the buffer of the last read rune. It is safe
// to call only once for a single read operation else it panics.
func (lx *lexer) unread() {
	if lx.lastRuneSize == 0 {
		panic("can't unread")
	}
	lx.pos -= lx.lastRuneSize
	lx.lastRuneSize = 0
}

// scan reads a sequence of characters that are in the set of allowed chars.
func (lx *lexer) scan(allowed []byte) {
	for lx.accept(allowed) {
	}
}

// accept reads next rune if it is in the set of allowed chars.
func (lx *lexer) accept(chars []byte) bool {
	r, err := lx.read()
	if err != nil {
		return false
	}
	if !bytes.ContainsRune(chars, r) {
		lx.unread()
		return false
	}
	return true
}
