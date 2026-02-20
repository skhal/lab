// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import "github.com/skhal/lab/x/fin/internal/fin"

// ReinvestDividend reinvests dividends.
func ReinvestDividend(pos fin.Position) fin.Position {
	pos.Investment += pos.Dividend
	pos.Dividend = 0
	return pos
}
