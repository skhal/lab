// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"fmt"

	dperrors "github.com/skhal/lab/dp/errors"
)

var errBufferFull = errors.New("buffer is full")

type testBuffer struct {
	buf []byte
	cap int
	len int
}

func newTestBuffer(cap int) *testBuffer {
	return &testBuffer{
		buf: make([]byte, cap),
		cap: cap,
	}
}

func (tb *testBuffer) Bytes() []byte {
	return tb.buf[:tb.len]
}

func (tb *testBuffer) Write(b []byte) (n int, err error) {
	n = min(len(b), tb.cap-tb.len)
	if n == 0 {
		return 0, errBufferFull
	}
	m := copy(tb.buf[tb.len:tb.cap], b[:n])
	tb.len += m
	return m, nil
}

func ExampleWriterWithError() {
	buf := newTestBuffer(16)
	w := dperrors.NewWriterWithError(buf)
	w.Write([]byte("hello"))
	w.Write([]byte(" "))
	w.Write([]byte("world"))
	w.Write([]byte(" from a test example"))
	w.Write([]byte(" from"))
	w.Write([]byte("a test example"))
	w.Write([]byte("!"))
	if err := w.Err(); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("buffer: %q\n", buf.Bytes())
	// Output:
	// error: buffer is full
	// buffer: "hello world from"
}
