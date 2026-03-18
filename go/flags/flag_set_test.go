// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags_test

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"testing"

	"github.com/skhal/lab/go/flags"
)

func ExampleFlagSet_noValidation() {
	args := []string{
		"-n", "-10",
		"-m", "-20",
	}
	var (
		n int
		m int
	)
	defer func() {
		fmt.Printf("n: %d\n", n)
		fmt.Printf("m: %d\n", m)
	}()
	fs := flags.NewFlagSet("example", flag.ExitOnError)
	fs.IntVar(&n, "n", 1, "positive number n")
	fs.IntVar(&m, "m", 2, "positive number m")
	if err := fs.ParseAndValidate(args); err != nil {
		if !errors.Is(err, flags.ErrFlag) {
			fmt.Println(err)
		}
		return
	}
	// Output:
	// n: -10
	// m: -20
}

func ExampleFlagSet_validate() {
	args := []string{
		"-n", "-10",
		"-m", "20",
	}
	var (
		n = 1
		m = 2
	)
	defer func() {
		fmt.Printf("n: %d\n", n)
		fmt.Printf("m: %d\n", m)
	}()
	fs := flags.NewFlagSet("example", flag.ExitOnError)
	fs.Var(newTestPositiveInt(&n), "n", "positive number n")
	fs.Var(newTestPositiveInt(&m), "m", "positive number m")
	if err := fs.ParseAndValidate(args); err != nil {
		if !errors.Is(err, flags.ErrFlag) {
			fmt.Println(err)
		}
		return
	}
	// Output:
	// n: -10
	// m: 20
}

func TestFlagSet_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fs      *flags.FlagSet
		args    []string
		wantErr error
	}{
		{
			name: "no flags",
			fs: func() *flags.FlagSet {
				t.Helper()
				return flags.NewFlagSet("test", flag.ContinueOnError)
			}(),
		},
		{
			name: "valid default value",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := 1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
		},
		{
			name: "invalid default value",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := -1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
			wantErr: flags.ErrFlag,
		},
		{
			name: "valid flag arg",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := -1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
			args: []string{"-n", "1"},
		},
		{
			name: "invalid flag arg",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := 1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
			args:    []string{"-n", "-1"},
			wantErr: flags.ErrFlag,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fs.Parse(tc.args); err != nil {
				t.Fatalf("Parse() unexpected error %v", err)
			}

			err := tc.fs.Validate()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate() unexpected error %v; want %v", err, tc.wantErr)
			}
		})
	}
}

func TestFlagSet_ParseAndValidate(t *testing.T) {
	tests := []struct {
		name    string
		fs      *flags.FlagSet
		args    []string
		wantErr error
	}{
		{
			name: "no flags",
			fs: func() *flags.FlagSet {
				t.Helper()
				return flags.NewFlagSet("test", flag.ContinueOnError)
			}(),
		},
		{
			name: "valid default value",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := 1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
		},
		{
			name: "invalid default value",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := -1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
			wantErr: flags.ErrFlag,
		},
		{
			name: "valid flag arg",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := -1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
			args: []string{"-n", "1"},
		},
		{
			name: "invalid flag arg",
			fs: func() *flags.FlagSet {
				t.Helper()
				fs := flags.NewFlagSet("test", flag.ContinueOnError)
				n := 1
				fs.Var(newTestPositiveInt(&n), "n", "test number")
				return fs
			}(),
			args:    []string{"-n", "-1"},
			wantErr: flags.ErrFlag,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fs.ParseAndValidate(tc.args)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate() unexpected error %v; want %v", err, tc.wantErr)
			}
		})
	}
}

type testPositiveInt struct {
	n *int
}

func newTestPositiveInt(n *int) flag.Value {
	return &testPositiveInt{n}
}

func (f testPositiveInt) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*f.n = n
	return nil
}

func (f testPositiveInt) String() string {
	if f.n == nil {
		return ""
	}
	return strconv.Itoa(*f.n)
}

func (f testPositiveInt) Validate() error {
	if *f.n <= 0 {
		return fmt.Errorf("non-positive integer")
	}
	return nil
}

type TestFlag struct {
	v flag.Value
	n string
	u string
}
