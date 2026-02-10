// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fin provides common units for market analysis.
package fin

import "fmt"

// Cents are the smallest unit of currency.
type Cents int64

// String implements fmt.Stringer interface.
func (c Cents) String() string {
	n := c / 100
	m := c % 100
	return fmt.Sprintf("%d.%02d", n, m)
}
