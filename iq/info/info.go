// Copyright 2025 Samvel Khalatyan. All rights reserved.

package info

import (
	"fmt"

	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
)

// Run prints questions from the registry.
func Run(reg *registry.R) error {
	reg.Visit(printQuestion)
	return nil
}

func printQuestion(q *pb.Question) {
	fmt.Printf("%d\t%s\n", q.GetId(), q.GetDescription())
}
