// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fin

import "time"

// Balance represents a state of the account as of date with cash available for
// investment and open positions.
type Balance struct {
	Date     time.Time // as of the first day of the month
	Cash     Cents     // not invested amount
	Position Position  // invested amount
}
