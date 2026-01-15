// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Lint-todo validates todo-comments.

Synopsis:

	check-todo file [file ...]
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

func main() {
	flag.Parse()
	cfg := todo.NewConfig()
	if err := todo.Run(cfg, flag.Args()...); err != nil {
		if !errors.Is(err, todo.ErrCheck) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
