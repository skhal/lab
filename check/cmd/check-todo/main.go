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
	"errors"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

func main() {
	cfg := todo.NewConfig()
	if err := todo.Run(cfg, os.Args[1:]...); err != nil {
		if !errors.Is(err, todo.ErrCheck) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
