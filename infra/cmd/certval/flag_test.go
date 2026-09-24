// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"testing"
)

func TestFlagError_Unwrap(t *testing.T) {
	want := errors.New("test error")

	got := errors.Unwrap(FlagError{want})

	if got != want {
		t.Errorf("Unwrap got %s; want %s", got, want)
	}
}

func TestParseFlags(t *testing.T) {
	errTest := errors.New("test error")
	tests := []struct {
		name      string
		validator flagValidatorFunc
		wantErr   error
	}{
		{
			name:      "no error validator",
			validator: func() error { return nil },
		},
		{
			name: "validator with error",
			validator: func() error {
				return errTest
			},
			wantErr: errTest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ParseFlags(tc.validator)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
		})
	}
}
