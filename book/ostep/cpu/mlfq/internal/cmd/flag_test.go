// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cmd"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
)

func TestPolicySpecFlag_Set(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    policy.Spec
		wantErr error
	}{
		{
			name: "zero value",
		},
		{
			name: "allotment",
			s:    "1:0:0",
			want: policy.Spec{Allotment: 1},
		},
		{
			name:    "allotment missing",
			s:       ":0:0",
			wantErr: cmd.ErrPolicySpec,
		},
		{
			name:    "allotment not a number",
			s:       "a:0:0",
			wantErr: cmd.ErrPolicySpec,
		},
		{
			name: "priorities",
			s:    "0:1:0",
			want: policy.Spec{Priorities: 1},
		},
		{
			name:    "priorities missing",
			s:       "0::0",
			wantErr: cmd.ErrPolicySpec,
		},
		{
			name:    "priorities not a number",
			s:       "0:a:0",
			wantErr: cmd.ErrPolicySpec,
		},
		{
			name: "boost cycles",
			s:    "0:0:1",
			want: policy.Spec{BoostCycles: 1},
		},
		{
			name:    "boost cycles missing",
			s:       "0:0:",
			wantErr: cmd.ErrPolicySpec,
		},
		{
			name:    "boost cycles not a number",
			s:       "0:0:a",
			wantErr: cmd.ErrPolicySpec,
		},
		{
			name: "non zero",
			s:    "1:2:3",
			want: policy.Spec{Allotment: 1, Priorities: 2, BoostCycles: 3},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got policy.Spec
			fg := cmd.NewPolicySpecFlag(&got)

			err := fg.Set(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Set(%q) = %v; want %v", tc.s, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Set(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}

func TestPolicySpecFlag_String(t *testing.T) {
	tests := []struct {
		name string
		spec policy.Spec
		want string
	}{
		{
			name: "zero value",
			want: "0:0:0",
		},
		{
			name: "allotment",
			spec: policy.Spec{Allotment: 1},
			want: "1:0:0",
		},
		{
			name: "priorities",
			spec: policy.Spec{Priorities: 1},
			want: "0:1:0",
		},
		{
			name: "boost cycles",
			spec: policy.Spec{BoostCycles: 1},
			want: "0:0:1",
		},
		{
			name: "non zero",
			spec: policy.Spec{Allotment: 1, Priorities: 2, BoostCycles: 3},
			want: "1:2:3",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fg := cmd.NewPolicySpecFlag(&tc.spec)

			got := fg.String()

			if tc.want != got {
				t.Errorf("String() = %q; want %q", got, tc.want)
				t.Logf("Spec: %v", tc.spec)
			}
		})
	}
}
