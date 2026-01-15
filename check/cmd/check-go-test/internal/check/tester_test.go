// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestEvent_Fail(t *testing.T) {
	tests := []struct {
		name  string
		event *check.Event
		want  bool
	}{
		{
			name:  "build event action output",
			event: withBuildEventAction(t, build.ActionOutput),
		},
		{
			name:  "build event action fail",
			event: withBuildEventAction(t, build.ActionFail),
			want:  true,
		},
		{
			name:  "test event action start",
			event: withTestEventAction(t, test.ActionStart),
		},
		{
			name:  "test event action run",
			event: withTestEventAction(t, test.ActionRun),
		},
		{
			name:  "test event action pause",
			event: withTestEventAction(t, test.ActionPause),
		},
		{
			name:  "test event action continue",
			event: withTestEventAction(t, test.ActionContinue),
		},
		{
			name:  "test event action pass",
			event: withTestEventAction(t, test.ActionPass),
		},
		{
			name:  "test event action benchmark",
			event: withTestEventAction(t, test.ActionBenchmark),
		},
		{
			name:  "test event action fail",
			event: withTestEventAction(t, test.ActionFail),
			want:  true,
		},
		{
			name:  "test event action output",
			event: withTestEventAction(t, test.ActionOutput),
		},
		{
			name:  "test event action skip",
			event: withTestEventAction(t, test.ActionSkip),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.event.Fail()
			if got != tc.want {
				t.Errorf("check.Event.Fail() got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestJSONUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		b       string
		want    *check.Event
		wantErr error
	}{
		{
			name:    "not json",
			b:       "test",
			wantErr: check.ErrNotJSON,
		},
		{
			name: "build event",
			b:    `{"ImportPath":"test","Action":"build-output"}`,
			want: withBuildEvent(t, &build.Event{ImportPath: "test", Action: build.ActionOutput}),
		},
		{
			name: "test event",
			b:    `{"Package":"test","Action":"output"}`,
			want: withTestEvent(t, &test.TestEvent{Package: "test", Action: test.ActionOutput}),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := check.JSONUnmarshal([]byte(tc.b))
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("check.JSONUnmarshal() unexpected error %q; want %q", err, tc.wantErr)
				t.Log(tc.b)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("check.JSONUnmarshal() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNewEventID(t *testing.T) {
	tests := []struct {
		name  string
		event *check.Event
		want  check.EventID
	}{
		{
			name:  "build event",
			event: withBuildEvent(t, &build.Event{ImportPath: "test"}),
			want:  check.EventID("test"),
		},
		{
			name:  "test event no package",
			event: withTestEvent(t, &test.TestEvent{Test: "test"}),
			want:  check.EventID("test"),
		},
		{
			name:  "test event with package",
			event: withTestEvent(t, &test.TestEvent{Package: "package", Test: "test"}),
			want:  check.EventID("package/test"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := check.NewEventID(tc.event)

			if tc.want != got {
				t.Errorf("test.NewEventID(%s) = %s; want %s", (*printableEvent)(tc.event), got, tc.want)
			}
		})
	}
}

func withBuildEventAction(t *testing.T, a build.Action) *check.Event {
	return &check.Event{
		BuildEvent: &build.Event{
			Action: a,
		},
	}
}

func withTestEventAction(t *testing.T, a test.Action) *check.Event {
	return &check.Event{
		TestEvent: &test.TestEvent{
			Action: a,
		},
	}
}

type printableEvent check.Event

func (e *printableEvent) String() string {
	if e.BuildEvent != nil {
		be := e.BuildEvent
		return fmt.Sprintf("{import-path:%q}", be.ImportPath)
	}
	if e.TestEvent != nil {
		te := e.TestEvent
		return fmt.Sprintf("{package:%q test:%q}", te.Package, te.Test)
	}
	return "{}"
}

func withBuildEvent(t *testing.T, e *build.Event) *check.Event {
	t.Helper()
	return &check.Event{BuildEvent: e}
}

func withTestEvent(t *testing.T, e *test.TestEvent) *check.Event {
	t.Helper()
	return &check.Event{TestEvent: e}
}
