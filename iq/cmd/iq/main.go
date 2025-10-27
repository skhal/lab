// Copyright 2025 Samvel Khalatyan. All rights reserved.

/*
Iq gives access to the interview questions registry.

Synopsis:

	iq

Without arguments, iq dumps a list of questions, sorted by ID.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
)

var file = flag.String("file", "iq/registry/questions.txtpb", "questions list (txtpb)")

func main() {
	flag.Parse()
	if err := run(*file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(file string) error {
	reg, err := registry.Load(file)
	if err != nil {
		return err
	}
	reg.Visit(printQuestion)
	return nil
}

func printQuestion(q *pb.Question) {
	fmt.Printf("%d\t%s\n", q.GetId(), q.GetDescription())
}
