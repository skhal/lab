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

func TestErrDuplicateQuestion_Is(t *testing.T) {
	err := new(registry.ErrDuplicateQuestion)

	got := errors.Is(err, registry.ErrRegistry)

	if want := true; got != want {
		t.Errorf("errors.Is(%T, registry.ErrRegistry) = %v; want %v", err, got, want)
	}
}

func newQuestion(t *testing.T, id int, desc string, tags ...string) *pb.Question {
	t.Helper()
	q := new(pb.Question)
	q.SetId(int32(id))
	q.SetDescription(desc)
	if len(tags) != 0 {
		q.SetTags(tags)
	}
	return q
}

func TestLoad(t *testing.T) {
	cfg := &registry.Config{File: "questions.txtpb"}
	_, err := registry.Load(cfg)

	if err != nil {
		t.Errorf("registry.Load(%v) = _, %s; want no error", cfg, err)
	}
}

func TestRegistry_WithQuestions_errorsOnDuplciates(t *testing.T) {
	hasQuestion := newQuestion(t, 1, "one")
	dupQuestion := newQuestion(t, 1, "two")
	qq := []*pb.Question{
		hasQuestion,
		dupQuestion,
	}

	_, err := registry.WithQuestions(qq)

	wantErr := &registry.ErrDuplicateQuestion{Has: hasQuestion, New: dupQuestion}
	if diff := cmp.Diff(wantErr, err, protocmp.Transform()); diff != "" {
		t.Errorf("registry.WithQuestions(%v) mismatch (-want, +got):\n%s", qq, diff)
	}
}

func TestRegistry_CreateQuestion(t *testing.T) {
	qq := []*pb.Question{
		newQuestion(t, 1, "one"),
	}
	reg, _ := registry.WithQuestions(qq)
	desc := "test-description"
	tags := []string{"tag-one", "tag-two"}

	got, err := reg.CreateQuestion(desc, tags)

	if err != nil {
		t.Fatalf("registry.R.CreateQuestion(%s, %v) = _, %s; want no error", desc, tags, err)
	}
	want := newQuestion(t, 2, desc, tags...)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("registry.R.CreateQuestion(%s, %v) mismatch (-want, +got):\n%s", desc, tags, diff)
	}
	if got, want := reg.Updated(), true; want != got {
		t.Errorf("registry.R.Updated() = %v; want %v", got, want)
	}
}

func TestRegistry_Get(t *testing.T) {
	tests := []struct {
		name string
		qq   []*pb.Question
		id   registry.QuestionID
		want *pb.Question
	}{
		{
			name: "empty",
			id:   registry.QuestionID(1),
		},
		{
			name: "hit",
			qq:   []*pb.Question{newQuestion(t, 1, "one")},
			id:   registry.QuestionID(1),
			want: newQuestion(t, 1, "one"),
		},
		{
			name: "miss",
			qq:   []*pb.Question{newQuestion(t, 1, "one")},
			id:   registry.QuestionID(2),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reg, err := registry.WithQuestions(tc.qq)
			if err != nil {
				t.Fatalf("registry.WithQuestions(%v) failed", tc.qq)
			}

			got := reg.Get(tc.id)

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("registry.Get(%d) mismatch (-want, +got):\n%s", tc.id, diff)
			}
		})
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
