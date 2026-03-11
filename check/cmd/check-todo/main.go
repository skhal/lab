// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-todo validates todo-comments.
//
// SYNOPSIS
//
//	check-todo file [file ...]
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

func main() {
	if err := todo.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
