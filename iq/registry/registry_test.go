// Copyright 2025 Samvel Khalatyan. All rights reserved.

package registry_test

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
	"google.golang.org/protobuf/testing/protocmp"
)

var update = flag.Bool("update", false, "update golden files")

func TestDuplicateQuestionError_Is(t *testing.T) {
	err := new(registry.DuplicateQuestionError)

	got := errors.Is(err, registry.ErrRegistry)

	if want := true; got != want {
		t.Errorf("errors.Is(%T, registry.ErrRegistry) = %v; want %v", err, got, want)
	}
}

func TestLoad(t *testing.T) {
	cfg := &registry.Config{File: "questions.txtpb"}
	_, err := registry.Load(cfg)

	if err != nil {
		t.Errorf("registry.Load(%v) = _, %s; want no error", cfg, err)
	}
}

func TestWrite(t *testing.T) {
	opts := []registry.Option{
		registry.QuestionOption(newQuestion(t, 1, "one")),
	}
	reg := mustCreateRegistry(t, opts...)
	tmpfile := filepath.Join(t.TempDir(), "registry.txtpb")
	cfg := &registry.Config{File: tmpfile}
	golden := tests.GoldenFile("testdata/registry_one_question.txtpb")

	err := registry.Write(reg, cfg)

	if err != nil {
		t.Fatalf("registry.Write() unexpected error: %v", err)
	}
	got := mustReadFile(t, tmpfile)
	if *update {
		golden.Write(t, got)
	}
	if diff := golden.Diff(t, got); diff != "" {
		t.Errorf("registry.Write() mismatch (-want, +got):\n%s", diff)
	}
}

func TestWrite_withHeader(t *testing.T) {
	opts := []registry.Option{
		registry.HeaderOption(strings.Split(`# proto-file: path/to/foo.proto
# proto-message: Foo`, "\n")),
		registry.QuestionOption(newQuestion(t, 1, "one")),
	}
	reg := mustCreateRegistry(t, opts...)
	tmpfile := filepath.Join(t.TempDir(), "registry.txtpb")
	cfg := &registry.Config{File: tmpfile}
	golden := tests.GoldenFile("testdata/registry_one_question_with_header.txtpb")

	err := registry.Write(reg, cfg)

	if err != nil {
		t.Fatalf("registry.Write() unexpected error: %v", err)
	}
	got := mustReadFile(t, tmpfile)
	if *update {
		golden.Write(t, got)
	}
	if diff := golden.Diff(t, got); diff != "" {
		t.Errorf("registry.Write() mismatch (-want, +got):\n%s", diff)
	}
}

func TestWrite_afterLoad(t *testing.T) {
	reg := mustLoad(t, "testdata/registry_one_question.txtpb")
	tmpfile := filepath.Join(t.TempDir(), "registry.txtpb")
	cfg := &registry.Config{File: tmpfile}
	golden := tests.GoldenFile("testdata/registry_one_question.txtpb")

	err := registry.Write(reg, cfg)

	if err != nil {
		t.Fatalf("registry.Write(_, %v) unexpected error: %v", cfg, err)
	}
	got := mustReadFile(t, tmpfile)
	// do not update golden
	if diff := golden.Diff(t, got); diff != "" {
		t.Errorf("registry.Write() mismatch (-want, +got):\n%s", diff)
	}
}

func TestWrite_afterLoadWithHeader(t *testing.T) {
	reg := mustLoad(t, "testdata/registry_one_question_with_header.txtpb")
	tmpfile := filepath.Join(t.TempDir(), "registry.txtpb")
	cfg := &registry.Config{File: tmpfile}
	golden := tests.GoldenFile("testdata/registry_one_question_with_header.txtpb")

	err := registry.Write(reg, cfg)

	if err != nil {
		t.Fatalf("registry.Write(_, %v) unexpected error: %v", cfg, err)
	}
	got := mustReadFile(t, tmpfile)
	// do not update golden
	if diff := golden.Diff(t, got); diff != "" {
		t.Errorf("registry.Write(_, %v) mismatch (-want, +got):\n%s", cfg, diff)
	}
}

func TestRegistry_With_errorsOnDuplciates(t *testing.T) {
	hasQuestion := newQuestion(t, 1, "one")
	dupQuestion := newQuestion(t, 1, "two")
	opts := []registry.Option{
		registry.QuestionOption(hasQuestion),
		registry.QuestionOption(dupQuestion),
	}

	_, err := registry.With(opts...)

	wantErr := &registry.DuplicateQuestionError{Has: hasQuestion, New: dupQuestion}
	if diff := cmp.Diff(wantErr, err, protocmp.Transform()); diff != "" {
		t.Errorf("registry.WithQuestions() error mismatch (-want, +got):\n%s", diff)
	}
}

func TestRegistry_CreateQuestion(t *testing.T) {
	opts := []registry.Option{
		registry.QuestionOption(newQuestion(t, 1, "one")),
	}
	reg := mustCreateRegistry(t, opts...)
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
			reg := mustCreateRegistry(t, registry.QuestionSetOption(tc.qq))

			got := reg.Get(tc.id)

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("registry.Get(%d) mismatch (-want, +got):\n%s", tc.id, diff)
			}
		})
	}
}

func TestRegistry_Visit_allOrderByID(t *testing.T) {
	opts := []registry.Option{
		registry.QuestionOption(newQuestion(t, 2, "two")),
		registry.QuestionOption(newQuestion(t, 1, "one")),
		registry.QuestionOption(newQuestion(t, 3, "three")),
	}
	reg := mustCreateRegistry(t, opts...)
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

func mustCreateRegistry(t *testing.T, opts ...registry.Option) *registry.R {
	t.Helper()
	reg, err := registry.With(opts...)
	if err != nil {
		t.Fatalf("registry.With() unexpected error: %v", err)
	}
	return reg
}

func mustReadFile(t *testing.T, file string) string {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read file %s: %v", file, err)
	}
	return string(data)
}

func mustLoad(t *testing.T, file string) *registry.R {
	t.Helper()
	cfg := &registry.Config{File: file}
	reg, err := registry.Load(cfg)
	if err != nil {
		t.Fatalf("load %s: %v", file, err)
	}
	return reg
}
