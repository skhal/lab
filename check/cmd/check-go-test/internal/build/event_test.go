// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
)

func TestActino_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		action  build.Action
		want    string
		wantErr error
	}{
		{
			name:   "build output",
			action: build.ActionOutput,
			want:   "build-output",
		},
		{
			name:   "build fail",
			action: build.ActionFail,
			want:   "build-fail",
		},
		{
			name:    "invalid action",
			action:  build.Action(123),
			wantErr: build.ErrInvalidAction,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.action)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("json.Marshal(%q) unexpected error %q; want %q", tc.action, err, tc.wantErr)
			}
			var want string
			if tc.want != "" {
				want = fmt.Sprintf("%q", tc.want)
			}
			if diff := cmp.Diff(want, string(got)); diff != "" {
				t.Errorf("json.Marshal(%q) mismatch (-want, +got):\n%s", tc.action, diff)
			}
		})
	}
}

func TestActino_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		b       string
		want    build.Action
		wantErr error
	}{
		{
			name: "build output",
			b:    "build-output",
			want: build.ActionOutput,
		},
		{
			name: "build fail",
			b:    "build-fail",
			want: build.ActionFail,
		},
		{
			name:    "invalid action",
			b:       "invalid-action",
			wantErr: build.ErrInvalidAction,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var got build.Action
			b := fmt.Sprintf("%q", tc.b)

			err := json.Unmarshal([]byte(b), &got)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("json.Unmarshal(%q) unexpected error %q; want %q", tc.b, err, tc.wantErr)
			}
			if tc.want != got {
				t.Errorf("json.Unmarshal(%q) got %s; want %s", tc.b, got, tc.want)
			}
		})
	}
}
