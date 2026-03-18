// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags_test

import (
	"errors"
	"flag"
	"strconv"
	"testing"

	"github.com/skhal/lab/go/flags"
)

func TestFlagError_Is(t *testing.T) {
	tests := []struct {
		name    string
		f       *flag.Flag
		err     error
		wantErr error
	}{
		{
			name:    "errflag",
			f:       &flag.Flag{Name: "test", Value: newTestValue(t, 123)},
			wantErr: flags.ErrFlag,
		},
		{
			name:    "errflag",
			f:       &flag.Flag{Name: "test", Value: newTestValue(t, 123)},
			err:     ErrTestA,
			wantErr: ErrTestA,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := flags.NewFlagError(tc.f, tc.err)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Is() unexpected match %v", tc.wantErr)
			}
		})
	}
}

func TestFlagError_Is_not(t *testing.T) {
	f := &flag.Flag{Name: "test", Value: newTestValue(t, 123)}
	err := flags.NewFlagError(f, ErrTestA)

	if doNotWant := ErrTestB; errors.Is(err, doNotWant) {
		t.Errorf("Is() unexpected match %v", doNotWant)
	}
}

func TestFlagError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		f       *flag.Flag
		err     error
		wantErr error
	}{
		{
			name: "errflag",
			f:    &flag.Flag{Name: "test", Value: newTestValue(t, 123)},
		},
		{
			name:    "errflag",
			f:       &flag.Flag{Name: "test", Value: newTestValue(t, 123)},
			err:     ErrTestA,
			wantErr: ErrTestA,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := flags.NewFlagError(tc.f, tc.err)

			got := errors.Unwrap(err)

			if got != tc.wantErr {
				t.Errorf("Unwrap() want %v", tc.wantErr)
			}
		})
	}
}

var (
	ErrTestA = errors.New("test error A")
	ErrTestB = errors.New("test error B")
)

func newTestValue(t *testing.T, n int) flag.Value {
	t.Helper()
	v := testValue(n)
	return &v
}

type testValue int

func (v *testValue) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*v = testValue(n)
	return nil
}

func (v testValue) String() string {
	return strconv.Itoa(int(v))
}
