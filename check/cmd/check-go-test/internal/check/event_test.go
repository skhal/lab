// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestEvent_Fail(t *testing.T) {
	tests := []struct {
		name  string
		event check.Event
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

func TestEvent_ID(t *testing.T) {
	tests := []struct {
		name  string
		event check.Event
		want  check.EventID
	}{
		{
			name:  "build event",
			event: (*check.BuildEvent)(&build.Event{ImportPath: "test"}),
			want:  check.EventID("test"),
		},
		{
			name:  "test event no package",
			event: (*check.TestEvent)(&test.TestEvent{Test: "test"}),
			want:  check.EventID("test"),
		},
		{
			name:  "test event with package",
			event: (*check.TestEvent)(&test.TestEvent{Package: "package", Test: "test"}),
			want:  check.EventID("package/test"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.event.ID()

			if tc.want != got {
				t.Errorf("Event.ID() = %s; want %s", got, tc.want)
				t.Logf("event:\n%s", formatEvent(tc.event))
			}
		})
	}
}

func withBuildEventAction(t *testing.T, a build.Action) check.Event {
	t.Helper()
	return (*check.BuildEvent)(&build.Event{Action: a})
}

func withTestEventAction(t *testing.T, a test.Action) check.Event {
	t.Helper()
	return (*check.TestEvent)(&test.TestEvent{Action: a})
}

func formatEvent(e check.Event) string {
	switch v := e.(type) {
	case *check.BuildEvent:
		return fmt.Sprintf("{import-path:%q}", v.ImportPath)
	case *check.TestEvent:
		return fmt.Sprintf("{package:%q test:%q}", v.Package, v.Test)
	}
	return "{}"
}
