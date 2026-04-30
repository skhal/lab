// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello demonstrates running Go in a browser using WebAssembly.
package main

import "fmt"

// Go code must be compiled as a program for WebAssembly to run it, i.e.
// it has to be placed inside a main() function.
func main() {
	fmt.Println("hello")
}
