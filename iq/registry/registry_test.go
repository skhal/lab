// Copyright 2025 Samvel Khalatyan. All rights reserved.

package registry_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestErrQuestion_Is(t *testing.T) {
	err := new(registry.ErrQuestion)

	got := errors.Is(err, registry.ErrRegistry)

	if want := true; got != want {
		t.Errorf("errors.Is(%T, registry.ErrRegistry) = %v; want %v", err, got, want)
	}
}

func newQuestion(t *testing.T, id int, desc string) *pb.Question {
	t.Helper()
	q := new(pb.Question)
	q.SetId(int32(id))
	q.SetDescription(desc)
	return q
}

func TestRegistry_WithQuestions_errorsOnDuplciates(t *testing.T) {
	dupQuestion := newQuestion(t, 1, "two")
	qq := []*pb.Question{
		newQuestion(t, 1, "one"),
		dupQuestion,
	}

	_, err := registry.WithQuestions(qq)

	wantErr := &registry.ErrQuestion{Question: dupQuestion}
	if diff := cmp.Diff(wantErr, err, protocmp.Transform()); diff != "" {
		t.Errorf("registry.WithQuestions(%v) mismatch (-want, +got):\n%s", qq, diff)
	}
}

func TestRegistry_Visit_allOrderByID(t *testing.T) {
	qq := []*pb.Question{
		newQuestion(t, 2, "two"),
		newQuestion(t, 1, "one"),
		newQuestion(t, 3, "three"),
	}
	reg, _ := registry.WithQuestions(qq)
	var got []int
	visitor := func(q *pb.Question) {
		got = append(got, int(q.GetId()))
	}

	reg.Visit(visitor)

	want := []int{1, 2, 3}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("registry.R.Visit() mismatch (-want, +got):\n%s", diff)
	}
}
