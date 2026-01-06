// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestAction_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		action test.Action
		want   string
	}{
		{
			name:   "unspecified",
			action: test.ActionUnspecified,
			want:   `"unspecified"`,
		},
		{
			name:   "start",
			action: test.ActionStart,
			want:   `"start"`,
		},
		{
			name:   "run",
			action: test.ActionRun,
			want:   `"run"`,
		},
		{
			name:   "pause",
			action: test.ActionPause,
			want:   `"pause"`,
		},
		{
			name:   "continue",
			action: test.ActionContinue,
			want:   `"cont"`,
		},
		{
			name:   "pass",
			action: test.ActionPass,
			want:   `"pass"`,
		},
		{
			name:   "benchmark",
			action: test.ActionBenchmark,
			want:   `"bench"`,
		},
		{
			name:   "fail",
			action: test.ActionFail,
			want:   `"fail"`,
		},
		{
			name:   "output",
			action: test.ActionOutput,
			want:   `"output"`,
		},
		{
			name:   "skip",
			action: test.ActionSkip,
			want:   `"skip"`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.action)

			if err != nil {
				t.Errorf("json.Marshal(%q) unexpected error %q", tc.action, err)
			}
			if diff := cmp.Diff(tc.want, string(b)); diff != "" {
				t.Errorf("json.Marshal(%q) mismatch (-want, +got):\n%s", tc.action, diff)
			}
		})
	}
}

func TestAction_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		b    string
		want test.Action
	}{
		{
			name: "unspecified",
			b:    `"unspecified"`,
			want: test.ActionUnspecified,
		},
		{
			name: "start",
			b:    `"start"`,
			want: test.ActionStart,
		},
		{
			name: "run",
			b:    `"run"`,
			want: test.ActionRun,
		},
		{
			name: "pause",
			b:    `"pause"`,
			want: test.ActionPause,
		},
		{
			name: "continue",
			b:    `"cont"`,
			want: test.ActionContinue,
		},
		{
			name: "pass",
			b:    `"pass"`,
			want: test.ActionPass,
		},
		{
			name: "benchmark",
			b:    `"bench"`,
			want: test.ActionBenchmark,
		},
		{
			name: "fail",
			b:    `"fail"`,
			want: test.ActionFail,
		},
		{
			name: "output",
			b:    `"output"`,
			want: test.ActionOutput,
		},
		{
			name: "skip",
			b:    `"skip"`,
			want: test.ActionSkip,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got test.Action
			err := json.Unmarshal([]byte(tc.b), &got)

			if err != nil {
				t.Errorf("json.Unmarshal(%q, _) unexpected error %q", tc.b, err)
			}
			if got != tc.want {
				t.Errorf("json.Unmarshal(%q, _) got %s; want %s", tc.b, got, tc.want)
			}
		})
	}
}

func TestEvent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		event    *test.TestEvent
		wantJSON string
		wantErr  error
	}{
		{
			name:     "zero value",
			event:    &test.TestEvent{},
			wantJSON: `{}`,
		},
		{
			name: "time",
			event: &test.TestEvent{
				Time: time.Date(2006, 01, 02, 03, 04, 05, 0, time.UTC),
			},
			wantJSON: `{"Time":"2006-01-02T03:04:05Z"}`,
		},
		{
			name: "action",
			event: &test.TestEvent{
				Action: test.ActionStart,
			},
			wantJSON: `{"Action":"start"}`,
		},
		{
			name: "package",
			event: &test.TestEvent{
				Package: "test-package",
			},
			wantJSON: `{"Package":"test-package"}`,
		},
		{
			name: "test",
			event: &test.TestEvent{
				Test: "test-name",
			},
			wantJSON: `{"Test":"test-name"}`,
		},
		{
			name: "elapsed",
			event: &test.TestEvent{
				Elapsed: time.Duration(1.23 * float64(time.Second)),
			},
			wantJSON: `{"Elapsed":1.23}`,
		},
		{
			name: "output",
			event: &test.TestEvent{
				Output: "test-output",
			},
			wantJSON: `{"Output":"test-output"}`,
		},
		{
			name: "failed build",
			event: &test.TestEvent{
				FailedBuild: "test-build",
			},
			wantJSON: `{"FailedBuild":"test-build"}`,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.event)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("json.Marshal() = %v; want %v", err, tc.wantErr)
				t.Logf("event:\n%+v", tc.event)
			}
			if diff := cmp.Diff(tc.wantJSON, string(data)); diff != "" {
				t.Errorf("json.Marshal() got unexpected data (-want, +got):\n%s", diff)
				t.Logf("event:\n%+v", tc.event)
			}
		})
	}
}

func TestTestEvent_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantEvent *test.TestEvent
		wantErr   error
	}{
		{
			name:      "zero value",
			json:      `{}`,
			wantEvent: &test.TestEvent{},
		},
		{
			name: "time",
			json: `{"Time":"2006-01-02T03:04:05Z"}`,
			wantEvent: &test.TestEvent{
				Time: time.Date(2006, 01, 02, 03, 04, 05, 0, time.UTC),
			},
		},
		{
			name: "action",
			json: `{"Action":"start"}`,
			wantEvent: &test.TestEvent{
				Action: test.ActionStart,
			},
		},
		{
			name: "package",
			json: `{"Package":"test-package"}`,
			wantEvent: &test.TestEvent{
				Package: "test-package",
			},
		},
		{
			name: "test",
			json: `{"Test":"test-name"}`,
			wantEvent: &test.TestEvent{
				Test: "test-name",
			},
		},
		{
			name: "elapsed",
			json: `{"Elapsed":1.23}`,
			wantEvent: &test.TestEvent{
				Elapsed: time.Duration(1.23 * float64(time.Second)),
			},
		},
		{
			name: "output",
			json: `{"Output":"test-output"}`,
			wantEvent: &test.TestEvent{
				Output: "test-output",
			},
		},
		{
			name: "failed build",
			json: `{"FailedBuild":"test-build"}`,
			wantEvent: &test.TestEvent{
				FailedBuild: "test-build",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var testEvent *test.TestEvent

			err := json.Unmarshal([]byte(tc.json), &testEvent)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("json.Unmarshal() = %v; want %v", err, tc.wantErr)
				t.Logf("json:\n%q", tc.json)
			}
			if diff := cmp.Diff(tc.wantEvent, testEvent); diff != "" {
				t.Errorf("json.Unmarshal() got unexpected event (-want, +got):\n%s", diff)
				t.Logf("json:\n%q", tc.json)
			}
		})
	}
}
