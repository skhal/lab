// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"testing"
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
)

// NewBalance testing helper creates a [fin.Balance] with year, month, cash,
// and optional positions.
func NewBalance(t *testing.T, y int, m time.Month, cash int64, pp ...fin.Position) fin.Balance {
	t.Helper()
	b := fin.Balance{
		Date: NewTime(t, y, m),
		Cash: fin.Cents(cash),
	}
	if len(pp) > 0 {
		b.Position = pp[0]
	}
	return b
}

// NewPosition testing helper creates a [fin.Position] with investments and
// dividend.
func NewPosition(t *testing.T, inv, div int64) fin.Position {
	t.Helper()
	return fin.Position{
		Investment: fin.Cents(inv),
		Dividend:   fin.Cents(div),
	}
}
