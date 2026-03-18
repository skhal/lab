// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags

import (
	"flag"
)

// Validator verifies that the flag has a valid value. It should return a
// shallow error, describing only the problem, e.g. "non-positive value".
type Validator interface {
	// Validate runs validation logic.
	Validate() error
}

// FlagSet add flags validation to [flag.FlagSet]. It calls *.Validate() on
// every flag that implements Validator interface:
//
//	type Validator interface {
//		Validate() error
//	}
//
// which should return a shallow error, describing only what is wrong. The
// [flags.FlagSet] adds context such as flag name and value to the error
// description.
//
// The FlagSet stops on the first invalid flag and returns FlagError.
//
// EXAMPLE
//
//		fs := flags.NewFlagSet("demo", flag.ExitOnError)
//	 // register flags
//		if err := fs.ParseAndValidate(); err != nil {
//			return err
//		}
type FlagSet struct {
	*flag.FlagSet
}

// NewFlagSet creates a FlagSet object that wraps [flag.FlagSet].
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	return &FlagSet{
		FlagSet: flag.NewFlagSet(name, errorHandling),
	}
}

// ParseAndValidate runs [flag.Parse], followed by [Validate]. It returns first
// non-nil error from either of these calls.
func (fs *FlagSet) ParseAndValidate(args []string) error {
	if err := fs.Parse(args); err != nil {
		return err
	}
	return fs.Validate()
}

// Validate runs flag validator for any flag that implements [Validator]
// interface.
func (fs *FlagSet) Validate() (err error) {
	vfn := func(f *flag.Flag) {
		if err != nil {
			return
		}
		v, ok := f.Value.(Validator)
		if !ok {
			return
		}
		if e := v.Validate(); e != nil {
			err = FlagError{f, e}
		}
	}
	// Validate all flags to check default values too.
	fs.VisitAll(vfn)
	return
}
