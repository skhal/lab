// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package registry_test

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/iq/cmd/iq/internal/pb"
	"github.com/skhal/lab/iq/cmd/iq/internal/registry"
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
	cfg := &registry.Config{File: "testdata/registry_one_question.txtpb"}
	_, err := registry.Load(cfg)

	if err != nil {
		t.Errorf("registry.Load(%v) = _, %s; want no error", cfg, err)
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name   string
		opts   []registry.Option
		golden tests.GoldenFile
	}{
		{
			name: "no header",
			opts: []registry.Option{
				registry.QuestionOption(newQuestion(t, 1, "one")),
			},
			golden: tests.GoldenFile("testdata/registry_one_question.txtpb"),
		},
		{
			name: "with header",
			opts: []registry.Option{
				registry.HeaderOption([]byte(`# proto-file: path/to/foo.proto
# proto-message: Foo`)),
				registry.QuestionOption(newQuestion(t, 1, "one")),
			},
			golden: tests.GoldenFile("testdata/registry_one_question_with_header.txtpb"),
		},
		{
			name: "root path",
			opts: []registry.Option{
				registry.RootPathOption("test/prefix"),
			},
			golden: tests.GoldenFile("testdata/registry_with_root_path.txtpb"),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reg := mustCreateRegistry(t, tc.opts...)
			tmpfile := filepath.Join(t.TempDir(), "registry.txtpb")
			cfg := &registry.Config{File: tmpfile}

			err := registry.Write(reg, cfg)

			if err != nil {
				t.Fatalf("registry.Write() unexpected error: %v", err)
			}
			got := mustReadFile(t, tmpfile)
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("registry.Write() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWrite_afterLoad(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		golden tests.GoldenFile
	}{
		{
			name:   "no header",
			file:   "testdata/registry_one_question.txtpb",
			golden: tests.GoldenFile("testdata/registry_one_question.txtpb"),
		},
		{
			name:   "with header",
			file:   "testdata/registry_one_question_with_header.txtpb",
			golden: tests.GoldenFile("testdata/registry_one_question_with_header.txtpb"),
		},
		{
			name:   "root path",
			file:   "testdata/registry_with_root_path.txtpb",
			golden: tests.GoldenFile("testdata/registry_with_root_path.txtpb"),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reg := mustLoad(t, tc.file)
			tmpfile := filepath.Join(t.TempDir(), "registry.txtpb")
			cfg := &registry.Config{File: tmpfile}

			err := registry.Write(reg, cfg)

			if err != nil {
				t.Fatalf("registry.Write(_, %v) unexpected error: %v", cfg, err)
			}
			got := mustReadFile(t, tmpfile)
			// do not update golden file
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("registry.Write() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestRegistry_With_errorsOnDuplicates(t *testing.T) {
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

func TestRegistry_GetByID(t *testing.T) {
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

			got := reg.GetByID(tc.id)

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("registry.GetByID(%d) mismatch (-want, +got):\n%s", tc.id, diff)
			}
		})
	}
}

func TestRegistry_GetByTag(t *testing.T) {
	tests := []struct {
		name string
		qq   []*pb.Question
		tag  registry.Tag
		want []*pb.Question
	}{
		{
			name: "empty",
			tag:  registry.Tag("foo"),
		},
		{
			name: "hit",
			qq:   []*pb.Question{newQuestion(t, 1, "one", "foo")},
			tag:  registry.Tag("foo"),
			want: []*pb.Question{newQuestion(t, 1, "one", "foo")},
		},
		{
			name: "hit several",
			qq: []*pb.Question{
				newQuestion(t, 1, "one", "foo"),
				newQuestion(t, 2, "two", "bar"),
				newQuestion(t, 3, "three", "foo"),
			},
			tag: registry.Tag("foo"),
			want: []*pb.Question{
				newQuestion(t, 1, "one", "foo"),
				newQuestion(t, 3, "three", "foo"),
			},
		},
		{
			name: "miss",
			qq:   []*pb.Question{newQuestion(t, 1, "one", "foo")},
			tag:  registry.Tag("bar"),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reg := mustCreateRegistry(t, registry.QuestionSetOption(tc.qq))

			got := reg.GetByTag(tc.tag)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Errorf("registry.GetByTag(%s) mismatch (-want, +got):\n%s", tc.tag, diff)
			}
		})
	}
}

func TestRegistry_GetTags(t *testing.T) {
	tests := []struct {
		name string
		qq   []*pb.Question
		want []registry.Tag
	}{
		{
			name: "empty",
		},
		{
			name: "one question one tag",
			qq:   []*pb.Question{newQuestion(t, 1, "one", "foo")},
			want: []registry.Tag{"foo"},
		},
		{
			name: "one question two tags",
			qq:   []*pb.Question{newQuestion(t, 1, "one", "foo", "bar")},
			want: []registry.Tag{"foo", "bar"},
		},
		{
			name: "two questions one tag",
			qq: []*pb.Question{
				newQuestion(t, 1, "one", "foo"),
				newQuestion(t, 2, "two", "foo"),
			},
			want: []registry.Tag{"foo"},
		},
		{
			name: "two questions two tags",
			qq: []*pb.Question{
				newQuestion(t, 1, "one", "foo"),
				newQuestion(t, 2, "two", "bar"),
			},
			want: []registry.Tag{"foo", "bar"},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reg := mustCreateRegistry(t, registry.QuestionSetOption(tc.qq))

			got := reg.GetTags()

			opts := []cmp.Option{
				cmp.Transformer("SortTags", func(tags []registry.Tag) []registry.Tag {
					// Do not mutate the input - copy
					tt := append([]registry.Tag(nil), tags...)
					sort.Slice(tt, func(i, j int) bool {
						return strings.Compare(string(tt[i]), string(tt[j])) < 0
					})
					return tt
				}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Errorf("registry.GetTags() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestRegistry_RootPath_option(t *testing.T) {
	reg := mustCreateRegistry(t, registry.RootPathOption("test/prefix"))

	got := reg.RootPath()

	if want := "test/prefix"; got != want {
		t.Errorf("registry.(*R).RootPath() = %q; want %q", got, want)
	}
}

func TestRegistry_RootPath_load(t *testing.T) {
	reg := mustLoad(t, "testdata/registry_with_root_path.txtpb")

	got := reg.RootPath()

	if want := "test/prefix"; got != want {
		t.Errorf("registry.(*R).RootPath() = %q; want %q", got, want)
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
