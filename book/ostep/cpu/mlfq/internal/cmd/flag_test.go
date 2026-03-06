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
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
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
			want: "0",
		},
		{
			name: "allotment",
			spec: policy.Spec{Allotment: 1},
			want: "1",
		},
		{
			name: "priorities",
			spec: policy.Spec{Priorities: 1},
			want: "0:1",
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

func TestProcSpecListFlag_Set(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    []*proc.Spec
		wantErr error
	}{
		{
			name: "zero value",
		},
		{
			name: "non zero",
			s:    "1:2:3:4,5:6:7:8",
			want: []*proc.Spec{
				{
					Arrive:           1,
					CPUCycles:        2,
					IOAfterCPUCycles: 3,
					IOCycles:         4,
				},
				{
					Arrive:           5,
					CPUCycles:        6,
					IOAfterCPUCycles: 7,
					IOCycles:         8,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []*proc.Spec
			fg := cmd.NewProcSpecListFlag(&got)

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

func TestProcSpecListFlag_String(t *testing.T) {
	tests := []struct {
		name  string
		specs []*proc.Spec
		want  string
	}{
		{
			name: "zero value",
		},
		{
			name: "arrive",
			specs: []*proc.Spec{
				{Arrive: 1},
			},
			want: "1",
		},
		{
			name: "cpu cycles",
			specs: []*proc.Spec{
				{CPUCycles: 1},
			},
			want: "0:1",
		},
		{
			name: "io after cpu cycles",
			specs: []*proc.Spec{
				{IOAfterCPUCycles: 1},
			},
			want: "0:0:1",
		},
		{
			name: "io cycles",
			specs: []*proc.Spec{
				{IOCycles: 1},
			},
			want: "0:0:0:1",
		},
		{
			name: "multiple specs",
			specs: []*proc.Spec{
				{Arrive: 10, CPUCycles: 11, IOAfterCPUCycles: 12, IOCycles: 13},
				{Arrive: 20, CPUCycles: 21, IOAfterCPUCycles: 22, IOCycles: 23},
			},
			want: "10:11:12:13,20:21:22:23",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fg := cmd.NewProcSpecListFlag(&tc.specs)

			got := fg.String()

			if tc.want != got {
				t.Errorf("String() = %q; want %q", got, tc.want)
				t.Logf("Specs: %v", tc.specs)
			}
		})
	}
}

func TestProcSpecFlag_Set(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    proc.Spec
		wantErr error
	}{
		{
			name: "zero value",
		},
		{
			name: "arrive",
			s:    "1",
			want: proc.Spec{Arrive: 1},
		},
		{
			name:    "arrive missing",
			s:       ":1",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name:    "arrive not a number",
			s:       "a:1",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name: "cpu cycles",
			s:    "0:1:0",
			want: proc.Spec{CPUCycles: 1},
		},
		{
			name:    "cpu cycles missing",
			s:       "0::0",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name:    "cpu cycles not a number",
			s:       "0:a:0",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name: "io after cpu cycles",
			s:    "0:0:1",
			want: proc.Spec{IOAfterCPUCycles: 1},
		},
		{
			name:    "io after cpu cycles missing",
			s:       "0:0::0",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name:    "io after cpu cycles not a number",
			s:       "0:0:a:0",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name: "io cycles",
			s:    "0:0:0:1",
			want: proc.Spec{IOCycles: 1},
		},
		{
			name:    "io cycles missing",
			s:       "0:0:0:",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name:    "io cycles not a number",
			s:       "0:0:0:a",
			wantErr: cmd.ErrProcSpec,
		},
		{
			name: "non zero",
			s:    "1:2:3:4",
			want: proc.Spec{
				Arrive:           1,
				CPUCycles:        2,
				IOAfterCPUCycles: 3,
				IOCycles:         4,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got proc.Spec
			fg := cmd.NewProcSpecFlag(&got)

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

func TestProcSpecFlag_String(t *testing.T) {
	tests := []struct {
		name string
		spec proc.Spec
		want string
	}{
		{
			name: "zero value",
			want: "0",
		},
		{
			name: "arrive",
			spec: proc.Spec{
				Arrive: 1,
			},
			want: "1",
		},
		{
			name: "cpu cycles",
			spec: proc.Spec{
				CPUCycles: 1,
			},
			want: "0:1",
		},
		{
			name: "io after cpu cycles",
			spec: proc.Spec{
				IOAfterCPUCycles: 1,
			},
			want: "0:0:1",
		},
		{
			name: "io cycles",
			spec: proc.Spec{
				IOCycles: 1,
			},
			want: "0:0:0:1",
		},
		{
			name: "valid specs",
			spec: proc.Spec{
				Arrive: 1, CPUCycles: 2, IOAfterCPUCycles: 3, IOCycles: 4,
			},
			want: "1:2:3:4",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fg := cmd.NewProcSpecFlag(&tc.spec)

			got := fg.String()

			if tc.want != got {
				t.Errorf("String() = %q; want %q", got, tc.want)
				t.Logf("Specs: %v", tc.spec)
			}
		})
	}
}
