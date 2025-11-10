// Copyright 2025 Samvel Khalatyan. All rights reserved.

/*
Lint-todo validates todo-comments.

Synopsis:

	lint-todo file [file ...]
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/lint-todo/internal/todo"
)

func main() {
	flag.Parse()
	cfg := todo.NewConfig()
	if err := todo.Run(cfg, flag.Args()...); err != nil {
		if !errors.Is(err, todo.ErrLint) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
