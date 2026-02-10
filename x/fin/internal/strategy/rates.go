// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import "github.com/skhal/lab/x/fin/internal/pb"

// Rate represents generic rate of something, e.g. rate of return.
type Rate float64

const rateNoChange Rate = 1.

// SPRateOfReturn calculates S&P's rate of return:
//
//	SPComposite / SPComposite_prev_month
func SPRateOfReturn(prev, curr *pb.Record) Rate {
	psp := float64(prev.GetQuote().GetSpComposite().GetCents())
	csp := float64(curr.GetQuote().GetSpComposite().GetCents())
	if psp == 0 || csp == 0 {
		// Since S&P is never equal to 0, a zero value for either previous or
		// current record indicates a missing record. We can't calculate RoR in
		// this case and keep the investment intact.
		return rateNoChange
	}
	return Rate(csp / psp)
}

// DivRateOfReturn calculates the rate of return of dividends. It is equal to
// dividends paid for one unit of S&P Composite currency:
//
//	Dividend / SPComposite
func DivRateOfReturn(r *pb.Record) Rate {
	div := float64(r.GetQuote().GetDividend().GetCents())
	sp := float64(r.GetQuote().GetSpComposite().GetCents())
	if sp == 0 {
		return 0
	}
	return Rate(div / sp)
}
