// Copyright 2025 Samvel Khalatyan. All rights reserved.

package info_test

import (
	"fmt"

	"github.com/skhal/lab/iq/info"
	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
)

func ExampleRun() {
	opts := []registry.Option{
		registry.QuestionOption(newQuestion(1, "demo question one", "tag-foo")),
		registry.QuestionOption(newQuestion(2, "demo question two", "tag-foo")),
		registry.QuestionOption(newQuestion(3, "demo question three", "tag-foo", "tag-bar")),
	}
	reg, err := registry.With(opts...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := info.Run(reg); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// 1	demo question one
	// 2	demo question two
	// 3	demo question three
}

func newQuestion(id int, desc string, tags ...string) *pb.Question {
	q := new(pb.Question)
	q.SetId(int32(id))
	q.SetDescription(desc)
	q.SetTags(tags)
	return q
}
