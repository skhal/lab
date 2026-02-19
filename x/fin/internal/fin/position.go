// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fin

// Position represents a single investment.
type Position struct {
	Investment Cents // the amount invested.
	Dividend   Cents // dividends paid out.
}

// Total returns sum of the investment and dividends in the position.
func (p Position) Total() Cents {
	return p.Investment + p.Dividend
}
