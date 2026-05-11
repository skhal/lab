// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestTestError_Error(t *testing.T) {
	tests := []struct {
		name   string
		events []*check.TestEvent
		want   string
	}{
		{
			name: "empty",
		},
		{
			name: "not output action",
			events: []*check.TestEvent{
				newTestEvent(t, test.ActionSkip, "test output to be excluded"),
			},
		},
		{
			name: "output action",
			events: []*check.TestEvent{
				newTestEvent(t, test.ActionOutput, "test output"),
			},
			want: "test output",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := check.TestError(tc.events)

			got := err.Error()

			if d := cmp.Diff(tc.want, got); d != "" {
				t.Errorf("TestError.Error() mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestBuildError_Error(t *testing.T) {
	tests := []struct {
		name   string
		events []*check.BuildEvent
		want   string
	}{
		{
			name: "empty",
		},
		{
			name: "not output action",
			events: []*check.BuildEvent{
				newBuildEvent(t, build.ActionFail, "test output to be excluded"),
			},
		},
		{
			name: "output action",
			events: []*check.BuildEvent{
				newBuildEvent(t, build.ActionOutput, "test output"),
			},
			want: "test output",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := check.BuildError(tc.events)

			got := err.Error()

			if d := cmp.Diff(tc.want, got); d != "" {
				t.Errorf("BuildError.Error() mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func ExampleCoverageError() {
	err := check.CoverageError{
		Package: "test/package",
		Got:     check.Coverage(10),
		Want:    check.Coverage(20),
	}
	fmt.Println(err.Error())
	// Output:
	// === COVERAGE: test/package
	//     coverage: 10.0% of statements
	//     threshold: 20.0%
	// --- FAIL
}

func newTestEvent(t *testing.T, a test.Action, out string) *check.TestEvent {
	t.Helper()
	e := &test.Event{Action: a, Output: out}
	return &check.TestEvent{Event: e}
}

func newBuildEvent(t *testing.T, a build.Action, out string) *check.BuildEvent {
	t.Helper()
	return (*check.BuildEvent)(&build.Event{Action: a, Output: out})
}
