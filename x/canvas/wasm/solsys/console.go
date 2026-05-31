// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package main

import "syscall/js"

type windowConsole struct {
	wc js.Value
}

func console() *windowConsole {
	return &windowConsole{
		wc: js.Global().Get("console"),
	}
}

// Error prints an error message to the console.
func (c *windowConsole) Error(s string) {
	c.wc.Call("error", s)
}
