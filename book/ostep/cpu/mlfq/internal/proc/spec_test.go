// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/proc"
)

func TestSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spec    proc.Spec
		wantErr error
	}{
		{
			name:    "zero value",
			wantErr: proc.ErrSpecCPUCycles,
		},
		{
			name:    "negative arrive",
			spec:    proc.Spec{Arrive: -1},
			wantErr: proc.ErrSpecArrive,
		},
		{
			name:    "negative cpu cycles",
			spec:    proc.Spec{CPUCycles: -1},
			wantErr: proc.ErrSpecCPUCycles,
		},
		{
			name:    "zero cpu cycles",
			spec:    proc.Spec{CPUCycles: 0},
			wantErr: proc.ErrSpecCPUCycles,
		},
		{
			name:    "negative io after cpu cycles",
			spec:    proc.Spec{CPUCycles: 1, IOAfterCPUCycles: -1},
			wantErr: proc.ErrSpecIOAfterCPUCycles,
		},
		{
			name: "zero io after cpu cycles",
			spec: proc.Spec{CPUCycles: 1, IOAfterCPUCycles: 0},
		},
		{
			name:    "io enabled negative io cycles",
			spec:    proc.Spec{CPUCycles: 1, IOAfterCPUCycles: 1, IOCycles: -1},
			wantErr: proc.ErrSpecIOCycles,
		},
		{
			name:    "io enabled zero io cycles",
			spec:    proc.Spec{CPUCycles: 1, IOAfterCPUCycles: 1, IOCycles: 0},
			wantErr: proc.ErrSpecIOCycles,
		},
		{
			name: "io disabled negative io cycles",
			spec: proc.Spec{CPUCycles: 1, IOCycles: -1},
		},
		{
			name: "io disabled zero io cycles",
			spec: proc.Spec{CPUCycles: 1, IOCycles: 0},
		},
		{
			name: "valid",
			spec: proc.Spec{
				Arrive:           1,
				CPUCycles:        2,
				IOAfterCPUCycles: 3,
				IOCycles:         4,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.Validate()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate() = %v; want %v", err, tc.wantErr)
			}
		})
	}
}
