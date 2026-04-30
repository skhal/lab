// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

// Wasm runs spreadsheets in browser using WebAssembler.
//
// SYNOPSIS
//
//	make -C ./x	/sheet/cmd/wasm
//	go run ./x/serve/ ./x/sheet/cmd/wasm/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/skhal/lab/x/sheet/internal/sheet"
)

var errJSFunc = errors.New("js func error")

func main() {
	fmt.Println("hello from Go")
	js.Global().Set("run", newJSFunc(run))
	// enter infinite loop to let JS interact with Go Wasm
	select {}
}

func newJSFunc(f func(string) (string, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if n := len(args); n != 1 {
			err := fmt.Errorf("%w: got %d args, want 1", errJSFunc, n)
			fmt.Println(err)
			return err
		}
		input := args[0].String()
		out, err := f(input)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return out
	})
}

func run(s string) (string, error) {
	st := sheet.New(sheet.WithVMEngine())
	sc := bufio.NewScanner(strings.NewReader(s))
	ln := 1
	for sc.Scan() {
		line := sc.Text()
		toks := strings.SplitN(line, " ", 2)
		if len(toks) != 2 {
			return "", fmt.Errorf(`L%d: want "ID VALUE", got %q`, ln, line)
		}
		if err := st.Set(toks[0], toks[1]); err != nil {
			return "", fmt.Errorf(`L%d: %s`, ln, err)
		}
		ln++
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	if err := st.Calculate(); err != nil {
		return "", err
	}
	var out strings.Builder
	st.VisitAll(func(id, _ string, val float64) bool {
		fmt.Fprintf(&out, "%s %.2f\n", id, val)
		return true
	})
	return out.String(), nil
}
