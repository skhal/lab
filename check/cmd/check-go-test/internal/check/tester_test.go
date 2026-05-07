// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestJSONUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		b       string
		want    check.Event
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
			want: (*check.BuildEvent)(
				&build.Event{ImportPath: "test", Action: build.ActionOutput},
			),
		},
		{
			name: "test event",
			b:    `{"Package":"test","Action":"output"}`,
			want: (*check.TestEvent)(
				&test.TestEvent{Package: "test", Action: test.ActionOutput},
			),
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
