// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fin

import (
	"fmt"
	"time"
)

// Quote is a the balance at the beginning of a given month.
// A value `2006 Jan: $1.23` means $1.23 balance as of Jan 1, 2006.
type Quote struct {
	Date    time.Time // first day of the month
	Balance Cents     // balance on the date
}

// String implements fmt.Stringer interface.
func (s Quote) String() string {
	d := s.Date.Format("2006 Jan")
	return fmt.Sprintf("%s %s", d, s.Balance)
}
