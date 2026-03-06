// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/policy"
)

func TestSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spec    policy.Spec
		wantErr error
	}{
		{
			name:    "zero value",
			wantErr: policy.ErrAllotment,
		},
		{
			name:    "negative allotment",
			spec:    policy.Spec{Allotment: -1},
			wantErr: policy.ErrAllotment,
		},
		{
			name:    "zero allotment",
			spec:    policy.Spec{Allotment: 0},
			wantErr: policy.ErrAllotment,
		},
		{
			name:    "negative priorities",
			spec:    policy.Spec{Allotment: 1, Priorities: -1},
			wantErr: policy.ErrPriorities,
		},
		{
			name:    "zero priorities",
			spec:    policy.Spec{Allotment: 1, Priorities: 0},
			wantErr: policy.ErrPriorities,
		},
		{
			name:    "negative boost cycles",
			spec:    policy.Spec{Allotment: 1, Priorities: 1, BoostCycles: -1},
			wantErr: policy.ErrBoostCycles,
		},
		{
			name:    "zero boost cycles",
			spec:    policy.Spec{Allotment: 1, Priorities: 1, BoostCycles: 0},
			wantErr: policy.ErrBoostCycles,
		},
		{
			name: "valid",
			spec: policy.Spec{Allotment: 1, Priorities: 1, BoostCycles: 1},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.Validate()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("(*Spec).Validate() = %v; want %v", err, tc.wantErr)
			}
		})
	}
}
