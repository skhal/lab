// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/cpu"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/io"
)

func TestRunner_Run(t *testing.T) {
	tests := []struct {
		name    string
		dur     cpu.Cycle
		runN    int
		wantErr error
	}{
		{
			name:    "zero duration",
			wantErr: io.ErrCompleted,
		},
		{
			name: "duration one",
			dur:  1,
		},
		{
			name:    "duration one completed",
			dur:     1,
			runN:    1,
			wantErr: io.ErrCompleted,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := io.NewRunner(tc.dur)
			for range tc.runN {
				r.Run()
			}

			err := r.Run()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Run() = %v; want %v", err, tc.wantErr)
			}
		})
	}
}

func TestRunner_Done(t *testing.T) {
	tests := []struct {
		name string
		dur  cpu.Cycle
		runN int
		want bool
	}{
		{
			name: "zero duration",
			want: true,
		},
		{
			name: "duration one no run",
			dur:  1,
		},
		{
			name: "duration one run one",
			dur:  1,
			runN: 1,
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := io.NewRunner(tc.dur)
			for range tc.runN {
				r.Run()
			}

			got := r.Done()

			if tc.want != got {
				t.Errorf("Done() = %v; want %v", got, tc.want)
			}
		})
	}
}
