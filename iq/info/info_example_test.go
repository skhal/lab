// Copyright 2025 Samvel Khalatyan. All rights reserved.

package info_test

import (
	"fmt"

	"github.com/skhal/lab/iq/info"
	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
)

func ExampleRun_printAll() {
	reg := mustCreateRegistry()
	cfg := new(info.Config)
	if err := info.Run(cfg, reg); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// 1	demo question one
	// 2	demo question two
	// 3	demo question three
}

func ExampleRun_printByID() {
	reg := mustCreateRegistry()
	cfg := new(info.Config)
	ids := []string{"2"}
	if err := info.Run(cfg, reg, ids...); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// 2	demo question two
}

func ExampleRun_printByTag() {
	reg := mustCreateRegistry()
	cfg := &info.Config{
		Tag: "foo",
	}
	if err := info.Run(cfg, reg); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// 1	demo question one
	// 3	demo question three
}

func ExampleRun_printByTagAndID() {
	reg := mustCreateRegistry()
	cfg := &info.Config{
		Tag: "foo",
	}
	ids := []string{"3"}
	if err := info.Run(cfg, reg, ids...); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// 3	demo question three
}

func mustCreateRegistry() *registry.R {
	opts := []registry.Option{
		registry.QuestionOption(newQuestion(1, "demo question one", "foo")),
		registry.QuestionOption(newQuestion(2, "demo question two", "bar")),
		registry.QuestionOption(newQuestion(3, "demo question three", "foo", "bar")),
	}
	reg, err := registry.With(opts...)
	if err != nil {
		panic(err)
	}
	return reg
}

func newQuestion(id int, desc string, tags ...string) *pb.Question {
	q := new(pb.Question)
	q.SetId(int32(id))
	q.SetDescription(desc)
	q.SetTags(tags)
	return q
}
