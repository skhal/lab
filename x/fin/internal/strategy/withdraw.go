// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"errors"
	"fmt"
	"math"

	"github.com/skhal/lab/x/fin/internal/fin"
)

// ErrWithdrawRate means the withdrawal rate is invalid.
var ErrWithdrawRate = errors.New("invalid rate")

const monthsInYear = 12

// YearlyWithdrawer withdraws percent amount from the position every 12 months.
type YearlyWithdrawer struct {
	Rate float64 // percent to withdraw, must be in range [0, 1]

	months int
}

// Rebalance updates the position by withdrawing [YearlyWithdrawer.Rate] every
// 12 months.
func (yw *YearlyWithdrawer) Rebalance(pos fin.Position) fin.Position {
	yw.months += 1
	if yw.months%monthsInYear == 0 {
		pos = yw.withdraw(pos)
	}
	return pos
}

func (yw *YearlyWithdrawer) withdraw(pos fin.Position) fin.Position {
	if yw.Rate < 0 || yw.Rate > 1 {
		panic(fmt.Errorf("rate %f: %w", yw.Rate, ErrWithdrawRate))
	}
	withdraw := func(c fin.Cents) fin.Cents {
		return c - fin.Cents(math.Floor(yw.Rate*float64(c)))
	}
	pos.Investment = withdraw(pos.Investment)
	pos.Dividend = withdraw(pos.Dividend)
	return pos
}
