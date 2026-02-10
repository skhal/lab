// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report

import (
	"io"

	"github.com/skhal/lab/x/fin/internal/fin"
)

// StrategyInfo summarizes the results of a strategy. It includes strategy name,
// description, and starting/end quotes.
type StrategyInfo struct {
	Name        string    // strategy name
	Description string    // strategy description
	Start       fin.Quote // strategy starting quote
	End         fin.Quote // strategy result
}

// Strategy generates a report for a single strategy.
func Strategy(w io.Writer, info StrategyInfo) error {
	return tmpls.ExecuteTemplate(w, "strategy.txt", info)
}
