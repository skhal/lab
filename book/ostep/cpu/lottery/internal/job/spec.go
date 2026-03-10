// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"errors"
	"fmt"
)

var (
	// ErrSpecLength means the field [Spec.Length] has invalid value.
	ErrSpecLength = errors.New("invalid length")

	// ErrSpecTickets means the field [Spec.Tickets] has invalid value.
	ErrSpecTickets = errors.New("invalid tickets")
)

const (
	minLength  = 1
	minTickets = 1
)

// Spec holds job configuration parameters.
type Spec struct {
	Length  int // length of the job.
	Tickets int // number of tickets the job owns.
}

// String implements [fmt.Stringer] interface.
func (s Spec) String() string {
	return fmt.Sprintf("len:%d tks:%d", s.Length, s.Tickets)
}

// Validate checks the [Spec] fields to having a meaningful number, e.g.,
// length and tickets fields should be positive numbers.
func (s *Spec) Validate() error {
	if s.Length < minLength {
		return ErrSpecLength
	}
	if s.Tickets < minTickets {
		return ErrSpecTickets
	}
	return nil
}
