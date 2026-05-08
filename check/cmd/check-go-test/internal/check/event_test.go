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

func ExampleCoverage_String() {
	fmt.Println(check.Coverage(0))
	fmt.Println(check.Coverage(12))
	fmt.Println(check.Coverage(100))
	// Output:
	// 0.0%
	// 12.0%
	// 100.0%
}

func TestBuildEvent_Fail(t *testing.T) {
	tests := []struct {
		name  string
		event *check.BuildEvent
		want  bool
	}{
		{
			name:  "action output",
			event: withBuildEventAction(t, build.ActionOutput),
		},
		{
			name:  "action fail",
			event: withBuildEventAction(t, build.ActionFail),
			want:  true,
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

func TestBuildEvent_ID(t *testing.T) {
	tests := []struct {
		name  string
		event *check.BuildEvent
		want  check.EventID
	}{
		{
			name:  "import path",
			event: (*check.BuildEvent)(&build.Event{ImportPath: "test"}),
			want:  check.EventID("test"),
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

func TestNewTestEvent(t *testing.T) {
	tests := []struct {
		name    string
		event   *test.Event
		want    *check.TestEvent
		wantErr error
	}{
		{
			name:  "no coverage",
			event: &test.Event{Output: "test"},
			want:  mustTestEvent(t, &test.Event{Output: "test"}),
		},
		{
			name:  "zero coverage",
			event: &test.Event{Output: "test\ncoverage: 0.0% of statements\n"},
			want: &check.TestEvent{
				Event: &test.Event{
					Output: "test\ncoverage: 0.0% of statements\n",
				},
				Coverage: new(check.Coverage(0)),
			},
		},
		{
			name:  "non zero coverage",
			event: &test.Event{Output: "test\ncoverage: 12.3% of statements\n"},
			want: &check.TestEvent{
				Event: &test.Event{
					Output: "test\ncoverage: 12.3% of statements\n",
				},
				Coverage: new(check.Coverage(12.3)),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := check.NewTestEvent(tc.event)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error %v; want %v", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestTestEvent_Fail(t *testing.T) {
	tests := []struct {
		name  string
		event *check.TestEvent
		want  bool
	}{
		{
			name:  "action start",
			event: mustTestEvent(t, &test.Event{Action: test.ActionStart}),
		},
		{
			name:  "action run",
			event: mustTestEvent(t, &test.Event{Action: test.ActionRun}),
		},
		{
			name:  "action pause",
			event: mustTestEvent(t, &test.Event{Action: test.ActionPause}),
		},
		{
			name:  "action continue",
			event: mustTestEvent(t, &test.Event{Action: test.ActionContinue}),
		},
		{
			name:  "action pass",
			event: mustTestEvent(t, &test.Event{Action: test.ActionPass}),
		},
		{
			name:  "action benchmark",
			event: mustTestEvent(t, &test.Event{Action: test.ActionBenchmark}),
		},
		{
			name:  "action fail",
			event: mustTestEvent(t, &test.Event{Action: test.ActionFail}),
			want:  true,
		},
		{
			name:  "action output",
			event: mustTestEvent(t, &test.Event{Action: test.ActionOutput}),
		},
		{
			name:  "action skip",
			event: mustTestEvent(t, &test.Event{Action: test.ActionSkip}),
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

func TestTestEvent_ID(t *testing.T) {
	tests := []struct {
		name  string
		event *check.TestEvent
		want  check.EventID
	}{
		{
			name:  "no package",
			event: mustTestEvent(t, &test.Event{Test: "test"}),
			want:  check.EventID("test"),
		},
		{
			name:  "with package",
			event: &check.TestEvent{Event: &test.Event{Package: "package", Test: "test"}},
			want:  check.EventID("package/test"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.event.ID()

			if tc.want != got {
				t.Errorf("(*TestEvent).ID() = %s; want %s", got, tc.want)
				t.Logf("event:\n%s", formatEvent(tc.event))
			}
		})
	}
}

func withBuildEventAction(t *testing.T, a build.Action) *check.BuildEvent {
	t.Helper()
	return (*check.BuildEvent)(&build.Event{Action: a})
}

func mustTestEvent(t *testing.T, te *test.Event) *check.TestEvent {
	t.Helper()
	event, err := check.NewTestEvent(te)
	if err != nil {
		t.Fatal(err)
	}
	return event
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
